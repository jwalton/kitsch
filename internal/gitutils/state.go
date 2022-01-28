package gitutils

import (
	"io/fs"
	"strings"

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
	}

	return result
}
