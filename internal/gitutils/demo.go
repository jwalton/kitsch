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
	// HeadDescription is the name of the current branch if HEAD is not detached, or else
	// a hash or the the name of a tag.
	HeadDescription string `yaml:"headDescription"`
	// IsDetached is true if HEAD is detached.
	IsDetached bool `yaml:"isDetached"`
	// IsTag is true if HeadDescription is for a tag.
	IsTag bool `yaml:"isTag"`
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
// localRef is compared to remoteRef.
func (git DemoGit) GetAheadBehind(localRef string, remoteRef string) (ahead int, behind int, err error) {
	if git.IsDetached || git.CurrentBranchUpstream == "" {
		return 0, 0, nil
	}
	if !git.IsDetached && (localRef == "refs/heads/"+git.HeadDescription) && remoteRef == "refs/remotes/"+git.CurrentBranchUpstream {
		return git.Ahead, git.Behind, nil
	}
	return 0, 0, fmt.Errorf("Unknown")
}

// Head returns information about the current head.
func (git DemoGit) Head(maxTagsToSearch int) (head HeadInfo, err error) {
	var headDescription string
	if !git.IsDetached {
		headDescription = git.HeadDescription
	} else if len(git.HeadDescription) == 40 && regexp.MustCompile("^[0-9a-f]+$").MatchString(git.HeadDescription) {
		headDescription = "(" + git.HeadDescription[0:shortSHALength] + ")"
	} else {
		headDescription = "(" + git.HeadDescription + ")"
	}

	return HeadInfo{
		Description: headDescription,
		Detached:    git.IsDetached,
		Hash:        git.HeadDescription,
		IsTag:       git.IsTag,
	}, nil
}

// State returns the current state of the repository.
func (git DemoGit) State() RepositoryState {
	return RepositoryState{
		State: git.CurrentState,
		Step:  git.Step,
		Total: git.Total,
	}
}

// Stats returns status counters for the given git repo.
func (git DemoGit) Stats() (GitStats, error) {
	return git.CurrentStats, nil
}
