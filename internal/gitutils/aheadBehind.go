package gitutils

import (
	"fmt"
	"strings"
)

// GetUpstream returns the upstream of the current branch if one exists, or
// an empty string otherwise.
func (utils *GitUtils) GetUpstream(branch string) string {
	config, err := utils.localConfig()
	if err != nil {
		return ""
	}

	branchConfig := config.Branches[branch]
	if branchConfig == nil {
		return ""
	}

	if !strings.HasPrefix(branchConfig.Merge, "refs/heads/") {
		return ""
	}

	return branchConfig.Remote + "/" + branchConfig.Merge[11:]
}

// GetAheadBehind returns how many commits ahead and behind the given
// branch is compared to compareToBranch.  You can use `HEAD` for the branch name.
func (utils *GitUtils) GetAheadBehind(branch string, compareToBranch string) (ahead int, behind int, err error) {
	aheadBehind, err := utils.git("rev-list", "--left-right", "--count", branch+"..."+compareToBranch)

	if err != nil {
		return 0, 0, err
	}

	fmt.Sscanf(aheadBehind, "%d %d", &ahead, &behind)
	return ahead, behind, nil
}
