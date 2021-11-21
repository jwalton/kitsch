package gitutils

import (
	"fmt"
	"regexp"
)

// DemoGit is an instance of the Git interface which returns demo values.  This
// is useful for testing, and for running Kitsch in "demo mode".
type DemoGit struct {
	// RepoRootDirectory is the path to the root directory of the git repo.
	RepoRootDirectory string `yaml:"repoDir"`
	// Head is the name of the current branch if HEAD is not detached, or else
	// a hash or the the name of a tag.
	Head string `yaml:"head"`
	// IsDetached is true if HEAD is detached.
	IsDetached bool `yaml:"isDetached"`
	// CurrentBranchUpstream is the current upstream branch, or "" if none.
	CurrentBranchUpstream string `yaml:"currentBranchUpstream"`

	// CurrentState is the current state of this repo.
	CurrentState RepositoryStateType `yaml:"state"`
	// Step is the current step number if we are rebasing, 0 otherwise.
	Step string `yaml:"step"`
	// Total is the total number of steps to complete to finish the rebase, or 0
	// if not rebasing.
	Total string `yaml:"total"`

	// StashCount is the current number of stashes.
	StashCount int `yaml:"stashCount"`
	// Ahead is the number of commits ahead of the upstream branch.
	Ahead int `yaml:"ahead"`
	// Behind is the number of commits behind the upstream branch.
	Behind int `yaml:"behind"`

	// Stats for the current git repo.
	CurrentStats GitStats `yaml:"stats"`
}

// RepoRoot returns the root of the git repository.
func (git DemoGit) RepoRoot() string {
	return git.RepoRootDirectory
}

// GetStashCount returns the number of stashes.
func (git DemoGit) GetStashCount() (int, error) {
	return git.StashCount, nil
}

// GetUpstream returns the upstream of the current branch if one exists, or
// an empty string otherwise.
func (git DemoGit) GetUpstream(branch string) string {
	if !git.IsDetached {
		return git.CurrentBranchUpstream
	}
	return ""
}

// GetAheadBehind returns how many commits ahead and behind the given
// branch is compared to compareToBranch.  You can use `HEAD` for the branch name.
func (git DemoGit) GetAheadBehind(branch string, compareToBranch string) (ahead int, behind int, err error) {
	if !git.IsDetached && (branch == git.Head || branch == "HEAD") && compareToBranch == git.CurrentBranchUpstream {
		return git.Ahead, git.Behind, nil
	}
	return 0, 0, fmt.Errorf("Unknown")
}

// State returns the current state of the repository.
func (git DemoGit) State() RepositoryState {
	var headDescription string
	if !git.IsDetached {
		headDescription = git.Head
	} else if len(git.Head) == 40 && regexp.MustCompile("^[0-9a-f]+$").MatchString(git.Head) {
		headDescription = "(" + git.Head[0:shortSHALength] + ")"
	} else {
		headDescription = "(" + git.Head + ")"
	}

	return RepositoryState{
		HeadDescription: headDescription,
		IsDetached:      git.IsDetached,
		State:           git.CurrentState,
		Step:            git.Step,
		Total:           git.Total,
	}
}

// Stats returns status counters for the given git repo.
func (git DemoGit) Stats() (GitStats, error) {
	return git.CurrentStats, nil
}
