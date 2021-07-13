package gitutils

import (
	"fmt"
	"strings"
)

// GetUpstream returns the upstream of the current branch if one exists, or
// an empty string otherwise.
func (utils *GitUtils) GetUpstream(branch string) string {
	upstream, err := utils.git("for-each-ref", "--format=%(upstream:short)", "refs/heads/"+branch)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(upstream)
}

// GetAheadBehind returns how many commits ahead and behind the given
// branch is compared to compareToBranch.  You can use `HEAD` for the branch name.
func (utils *GitUtils) GetAheadBehind(branch string, compareToBranch string) (ahead int, behind int, err error) {
	aheadBehind, err := utils.git("rev-list", "--left-right", "--count", branch+"..."+compareToBranch)

	if err != nil {
		return 0, 0, err
	}

	fmt.Sscanf(aheadBehind, "%i %i", &ahead, &behind)
	return ahead, behind, nil
}
