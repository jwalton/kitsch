package gitutils

import (
	"path/filepath"
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
	// Base is the name of the branch we are in, or the tag name or short hash
	// if partially detached.
	Base string
	// Branch is the name of the branch we are in, or empty string if the head
	// is detached.
	Branch string
}

// State returns the current state of the repository.
// Based loosely on posh-git's Get-GitBranch.
// https://github.com/dahlbyk/posh-git/blob/b79c2dc39c9387847642bc3b38fa2186b29f6113/src/GitUtils.ps1#L62
func (utils *GitUtils) State() RepositoryState {
	var result RepositoryState
	gitFolder := filepath.Join(utils.RepoRoot, ".git")

	rebaseMerge := filepath.Join(gitFolder, "rebase-merge")
	if fileutils.FileExists(rebaseMerge) {
		if fileutils.FileExists(filepath.Join(rebaseMerge, "interactive")) {
			result.State = StateRebasingInteractive
		} else {
			result.State = StateRebaseMerging
		}

		result.Branch = extractBranchName(fileutils.ReadFile(filepath.Join(rebaseMerge, "head-name")))
		result.Base = result.Branch
		result.Step = fileutils.ReadFile(filepath.Join(rebaseMerge, "msgnum"))
		result.Total = fileutils.ReadFile(filepath.Join(rebaseMerge, "end"))
	} else {
		rebaseApply := filepath.Join(gitFolder, "rebase-apply")
		if fileutils.FileExists(rebaseApply) {
			result.Step = fileutils.ReadFile(filepath.Join(rebaseApply, "next"))
			result.Total = fileutils.ReadFile(filepath.Join(rebaseApply, "last"))

			if fileutils.FileExists(filepath.Join(rebaseApply, "rebasing")) {
				result.State = StateRebasing
			} else if fileutils.FileExists(filepath.Join(rebaseApply, "applying")) {
				result.State = StateAMing
			} else {
				result.State = StateRebaseAMing
			}
		} else if fileutils.FileExists(filepath.Join(gitFolder, "MERGE_HEAD")) {
			result.State = StateMerging
		} else if fileutils.FileExists(filepath.Join(gitFolder, "CHERRY_PICK_HEAD")) {
			result.State = StateCherryPicking
		} else if fileutils.FileExists(filepath.Join(gitFolder, "REVERT_HEAD")) {
			result.State = StateReverting
		} else if fileutils.FileExists(filepath.Join(gitFolder, "BISECT_LOG")) {
			result.State = StateBisecting
		} else {
			result.State = StateNone
		}

		result.Base, result.Branch = utils.getBranchName()
	}

	return result
}

func (utils *GitUtils) getBranchName() (base string, branch string) {
	branch, err := utils.git("symbolic-ref", "HEAD")
	branch = extractBranchName(branch)
	base = branch

	// If that fails, get tag name
	if err != nil || base == "" {
		base, err = utils.git("describe", "--tags", "--exact-match HEAD")
		if base != "" {
			base = "(" + strings.TrimSpace(base) + ")"
		}
	}

	// If that fails, get the SHA
	if err != nil || base == "" {
		base, err = utils.git("rev-parse", "HEAD")
		if base != "" && len(base) > shortSHALength {
			base = "(" + strings.TrimSpace(base)[0:shortSHALength] + "...)"
		}
	}

	if err != nil {
		base = "???"
	}

	return base, branch
}

func extractBranchName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.TrimPrefix(name, "refs/heads/")
	return name
}
