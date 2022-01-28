package gitutils

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestStateOnBranch(t *testing.T) {
	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("ref: refs/heads/master\n"),
		},
		".git/refs/heads/master": &fstest.MapFile{
			Data: []byte("7c088a39dcd2dcda89f4dee1fd3eb41c1d34ea2f\n"),
		},
	}

	git := testGitUtils("/Users/oriana/dev/kitsch", files)

	state := git.State()
	assert.Equal(t,
		RepositoryState{
			State: StateNone,
			Step:  "",
			Total: "",
		},
		state,
	)
}
