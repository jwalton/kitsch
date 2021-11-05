package gitutils

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
)

const shortSHALength = 7

// RepositoryStateType represents the current state of a repository (e.g. rebasing, merging, etc...)
type RepositoryStateType string

const (
	// StateNone is the default state of a git repo.
	StateNone RepositoryStateType = ""
	// StateRebasingInteractive is for a git repo in the middle of an interactive rebase.
	StateRebasingInteractive RepositoryStateType = "REBASE-i"
	// StateRebaseMerging is for a git repo in the middle of a merge.
	StateRebaseMerging RepositoryStateType = "REBASE-m"
	// StateRebasing is for a git repo in the middle of a rebase operation.
	StateRebasing RepositoryStateType = "REBASE"
	// StateAMing is for a git repo in the middle of an `am` operation.
	StateAMing RepositoryStateType = "AM"
	// StateRebaseAMing is for a git repo where we don't know if the repo is rebasing or aming.
	StateRebaseAMing RepositoryStateType = "REBASE/AM"
	// StateMerging is for a git repo in the middle of a merge.
	StateMerging RepositoryStateType = "MERGING"
	// StateCherryPicking is for a git repo in the middle of a cherry-pick.
	StateCherryPicking RepositoryStateType = "CHERRY-PICKING"
	// StateReverting is for a git repo in the middle of a revert.
	StateReverting RepositoryStateType = "REVERTING"
	// StateBisecting is for a git repo in the middle of a bisect.
	StateBisecting RepositoryStateType = "BISECTING"
)

// RepositoryState represents the overall state of a git repository.
type RepositoryState struct {
	// State is the current state of this repo.
	State RepositoryStateType
	// Step is the current step number if we are rebasing, 0 otherwise.
	Step string
	// Total is the total number of steps to complete to finish the rebase, or 0
	// if not rebasing.
	Total string
	// HeadDescription is the name of the branch we are currently on if the head
	// is not detached.  If the head is detached, this will be  the branch name
	// if we are in the middle of a rebase or merge, the tag name if the head is
	// at a tag, or the short hash otherwise.
	HeadDescription string
	// IsDetached is true if the head is currently detached.
	IsDetached bool
}

func (g *GitUtils) readFileIfExist(path string) string {
	contents, err := fs.ReadFile(g.files, path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(contents))
}

// State returns the current state of the repository.
// Based loosely on posh-git's Get-GitBranch.
// https://github.com/dahlbyk/posh-git/blob/b79c2dc39c9387847642bc3b38fa2186b29f6113/src/GitUtils.ps1#L62
func (g *GitUtils) State() RepositoryState {
	var result RepositoryState

	if fileutils.FSFileExists(g.files, ".git/rebase-merge") {
		if fileutils.FSFileExists(g.files, ".git/rebase-merge/interactive") {
			result.State = StateRebasingInteractive
		} else {
			result.State = StateRebaseMerging
		}

		result.HeadDescription = extractBranchName(g.readFileIfExist(".git/rebase-merge/head-name"))
		result.IsDetached = true
		result.Step = g.readFileIfExist(".git/rebase-merge/msgnum")
		result.Total = g.readFileIfExist(".git/rebase-merge/end")
	} else {
		if fileutils.FSFileExists(g.files, ".git/rebase-apply") {
			result.Step = g.readFileIfExist(".git/rebase-apply/next")
			result.Total = g.readFileIfExist(".git/rebase-apply/last")

			if fileutils.FSFileExists(g.files, ".git/rebase-apply/rebasing") {
				result.State = StateRebasing
			} else if fileutils.FSFileExists(g.files, ".git/rebase-apply/applying") {
				result.State = StateAMing
			} else {
				result.State = StateRebaseAMing
			}
		} else if fileutils.FSFileExists(g.files, ".git/MERGE_HEAD") {
			result.State = StateMerging
		} else if fileutils.FSFileExists(g.files, ".git/CHERRY_PICK_HEAD") {
			result.State = StateCherryPicking
		} else if fileutils.FSFileExists(g.files, ".git/REVERT_HEAD") {
			result.State = StateReverting
		} else if fileutils.FSFileExists(g.files, ".git/BISECT_LOG") {
			result.State = StateBisecting
		} else {
			result.State = StateNone
		}

		result.HeadDescription, result.IsDetached = g.getHeadDescription()
	}

	return result
}

// getHeadDescription returns a description of the current "HEAD".  This will
// be the branch name, or if the head is detached this will be the tag name or
// short hash.
func (g *GitUtils) getHeadDescription() (description string, isDetached bool) {
	isDetached = true

	head := g.readFileIfExist(".git/HEAD")
	if head == "" {
		return "???", true
	}

	if strings.HasPrefix(head, "ref: ") {
		// The HEAD is a symbolic reference.
		if strings.HasPrefix(head[5:], "refs/heads/") {
			// If the HEAD file is a symbolic reference to a branch, extract the branch name.
			description = head[16:]
			isDetached = false
		} else {
			// If the HEAD is a symbolic reference to a non-branch (can this
			// even happen??), resolve it to a hash.
			var err error
			head, err = g.resolveSymbolicRef(head)
			if err != nil {
				return "???", true
			}
		}
	}

	// If we don't have a description, try to get a tag name
	if description == "" {
		tag, err := g.GetTagNameForHash(head)
		if err == nil && tag != "" {
			description = "(" + strings.TrimSpace(tag) + ")"
		}
	}

	// If that fails, use the hash
	if description == "" && len(head) > shortSHALength {
		description = "(" + head[0:shortSHALength] + "...)"
	}

	if description == "" {
		description = "???"
	}

	return description, isDetached
}

var errNotFound = errors.New("Not found")

// GetTagNameForHash returns the tag name for the hash, or an error if no such
// tag exists.  "hash" can be a short hash.
func (g *GitUtils) GetTagNameForHash(hash string) (string, error) {
	tagFiles, err := fs.ReadDir(g.files, ".git/refs/tags")
	if err != nil {
		return "", err
	}

	for _, tagFile := range tagFiles {
		content, err := fs.ReadFile(g.files, ".git/refs/tags/"+tagFile.Name())
		if err != nil {
			continue
		}
		if len(content) >= len(hash) && string(content[0:len(hash)]) == hash {
			return tagFile.Name(), nil
		}
	}

	return "", errNotFound
}

func (g *GitUtils) resolveSymbolicRef(ref string) (string, error) {
	// Resolve the symbolic ref to a hash.
	hashBytes, err := fs.ReadFile(g.files, ".git/"+ref[5:])
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(hashBytes)), nil
}

func extractBranchName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.TrimPrefix(name, "refs/heads/")
	return name
}
