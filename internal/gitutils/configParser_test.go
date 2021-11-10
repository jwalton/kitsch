package gitutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigParser(t *testing.T) {
	config := `
[core]
	repositoryformatversion = 0
[branch "master"]
	remote = origin
	merge = refs/heads/master
[branch "feature/projects"]
	remote = origin
	merge = refs/heads/feature/projects
`

	expected := gitconfig{
		Branches: map[string]*gitconfigBranch{
			"master": {
				Branch: "master",
				Remote: "origin",
				Merge:  "refs/heads/master",
			},
			"feature/projects": {
				Branch: "feature/projects",
				Remote: "origin",
				Merge:  "refs/heads/feature/projects",
			},
		},
	}

	result, err := parseGitConfig(strings.NewReader(config))
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
