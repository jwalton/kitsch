package gitutils

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jwalton/kitsch/internal/fileutils"
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
	State RepositoryStateType `yaml:"state"`
	// Step is the current step number if we are rebasing, 0 otherwise.
	Step string `yaml:"step"`
	// Total is the total number of steps to complete to finish the rebase, or 0
	// if not rebasing.
	Total string `yaml:"total"`
	// HeadDescription is the name of the branch we are currently on if the head
	// is not detached.  If the head is detached, this will be  the branch name
	// if we are in the middle of a rebase or merge, the tag name if the head is
	// at a tag, or the short hash otherwise.
	HeadDescription string `yaml:"headDescription"`
	// IsDetached is true if the head is currently detached.
	IsDetached bool `yaml:"isDetached"`
}

func (g *gitUtils) readFileIfExist(path string) string {
	if g.fsys == nil {
		return ""
	}

	contents, err := fs.ReadFile(g.fsys, path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(contents))
}

// State returns the current state of the repository.
// Based loosely on posh-git's Get-GitBranch.
// https://github.com/dahlbyk/posh-git/blob/b79c2dc39c9387847642bc3b38fa2186b29f6113/src/GitUtils.ps1#L62
func (g *gitUtils) State() RepositoryState {
	var result RepositoryState

	if g.fsys == nil {
		return result
	}

	if fileutils.FSFileExists(g.fsys, ".git/rebase-merge") {
		if fileutils.FSFileExists(g.fsys, ".git/rebase-merge/interactive") {
			result.State = StateRebasingInteractive
		} else {
			result.State = StateRebaseMerging
		}

		result.HeadDescription = extractBranchName(g.readFileIfExist(".git/rebase-merge/head-name"))
		result.IsDetached = true
		result.Step = g.readFileIfExist(".git/rebase-merge/msgnum")
		result.Total = g.readFileIfExist(".git/rebase-merge/end")
	} else {
		if fileutils.FSFileExists(g.fsys, ".git/rebase-apply") {
			result.Step = g.readFileIfExist(".git/rebase-apply/next")
			result.Total = g.readFileIfExist(".git/rebase-apply/last")

			if fileutils.FSFileExists(g.fsys, ".git/rebase-apply/rebasing") {
				result.State = StateRebasing
			} else if fileutils.FSFileExists(g.fsys, ".git/rebase-apply/applying") {
				result.State = StateAMing
			} else {
				result.State = StateRebaseAMing
			}
		} else if fileutils.FSFileExists(g.fsys, ".git/MERGE_HEAD") {
			result.State = StateMerging
		} else if fileutils.FSFileExists(g.fsys, ".git/CHERRY_PICK_HEAD") {
			result.State = StateCherryPicking
		} else if fileutils.FSFileExists(g.fsys, ".git/REVERT_HEAD") {
			result.State = StateReverting
		} else if fileutils.FSFileExists(g.fsys, ".git/BISECT_LOG") {
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
func (g *gitUtils) getHeadDescription() (description string, isDetached bool) {
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
func (g *gitUtils) GetTagNameForHash(hash string) (string, error) {
	// Check lightweight tags
	tags, err := g.repo.Tags()
	if err == nil {
		for ref, err := tags.Next(); err == nil && ref != nil; ref, err = tags.Next() {
			if g.hashMatchesTag(hash, ref.Hash()) {
				result := strings.TrimPrefix(string(ref.Name()), "refs/tags/")
				return result, nil
			}

		}
	}

	return "", errNotFound
}

// hashMatchesTag returns true if the `hash` represents the same commit as
// the given tag.  `hash` can be a "short hash".  This will return true
// if hash is a prefix of tagHash (the "lightweight" tag case) or if
// the tagHash is an annotated hash and hash is the hash of the tag.
func (g *gitUtils) hashMatchesTag(hash string, tagHash plumbing.Hash) bool {
	// If the tag is a lightweight tag, the hash should match.
	if strings.HasPrefix(tagHash.String(), hash) {
		return true
	}

	// Check to see if this is an annotated tag
	tagObj, err := g.repo.TagObject(tagHash)
	if err == nil {
		if strings.HasPrefix(tagObj.Target.String(), hash) {
			return true
		}
	}

	return false
}

// resolveSymbolicRef returns the hash for a given symbolic ref.
// e.g. this turns "refs/heads/master" into a hash.
func (g *gitUtils) resolveSymbolicRef(ref string) (string, error) {
	hash, err := g.repo.ResolveRevision(plumbing.Revision(ref))
	return hash.String(), err
}

// extractBranchName returns the branch name from a symbolic ref, or returns
// the passed in string otherwise.
func extractBranchName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.TrimPrefix(name, "refs/heads/")
	return name
}
