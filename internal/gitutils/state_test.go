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

	git := &gitUtils{
		pathToGit: "git",
		fsys:      files,
		repoRoot:  "/Users/oriana/dev/kitsch-prompt",
	}

	state := git.State()
	assert.Equal(t,
		RepositoryState{
			State:           StateNone,
			Step:            "",
			Total:           "",
			HeadDescription: "master",
			IsDetached:      false,
		},
		state,
	)
}

func TestStateDetached(t *testing.T) {
	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("0123456789abcdef0123456789abcdef01234567\n"),
		},
		".git/refs/heads/master": &fstest.MapFile{
			Data: []byte("7c088a39dcd2dcda89f4dee1fd3eb41c1d34ea2f\n"),
		},
	}

	git := &gitUtils{
		pathToGit: "git",
		fsys:      files,
		repoRoot:  "/Users/oriana/dev/kitsch-prompt",
	}

	state := git.State()
	assert.Equal(t,
		RepositoryState{
			State:           StateNone,
			Step:            "",
			Total:           "",
			HeadDescription: "(0123456...)",
			IsDetached:      true,
		},
		state,
	)
}

func TestStateOnTag(t *testing.T) {
	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("0123456789abcdef0123456789abcdef01234567\n"),
		},
		".git/refs/heads/master": &fstest.MapFile{
			Data: []byte("7c088a39dcd2dcda89f4dee1fd3eb41c1d34ea2f\n"),
		},
		".git/refs/tags/v1.0.0": &fstest.MapFile{
			Data: []byte("0123456789abcdef0123456789abcdef01234567\n"),
		},
	}

	git := &gitUtils{
		pathToGit: "git",
		fsys:      files,
		repoRoot:  "/Users/oriana/dev/kitsch-prompt",
	}

	state := git.State()
	assert.Equal(t,
		RepositoryState{
			State:           StateNone,
			Step:            "",
			Total:           "",
			HeadDescription: "(v1.0.0)",
			IsDetached:      true,
		},
		state,
	)
}

func TestStateInNonGitRepo(t *testing.T) {
	git := &gitUtils{
		pathToGit: "git",
		fsys:      nil,
		repoRoot:  "",
	}

	state := git.State()
	assert.Equal(t,
		RepositoryState{
			State:           StateNone,
			Step:            "",
			Total:           "",
			HeadDescription: "",
			IsDetached:      false,
		},
		state,
	)
}
