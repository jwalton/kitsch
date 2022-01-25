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

	git := testGitUtils("/Users/oriana/dev/kitsch", files)

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

	git := testGitUtils("/Users/oriana/dev/kitsch", files)

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

func TestStateOnAnnotatedTag(t *testing.T) {
	files := fstest.MapFS{
		// Head points at commit aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
		".git/HEAD": &fstest.MapFile{
			Data: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n"),
		},
		".git/refs/heads/master": &fstest.MapFile{
			Data: []byte("7c088a39dcd2dcda89f4dee1fd3eb41c1d34ea2f\n"),
		},
		// Annotated tag points to object at bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
		".git/refs/tags/v1.1.0": &fstest.MapFile{
			Data: []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\n"),
		},
		// The tag object points to the commit.
		".git/objects/bb/bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb": &fstest.MapFile{
			Data: generateGitObject("tag", "object aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\ntype commit\ntag v1.1.0\ntagger Jason Walton <dev@lucid.thedreaming.org> 1642726592 -0500\n\nv1.1.0\n"),
		},
	}

	git := testGitUtils("/Users/oriana/dev/kitsch", files)

	state := git.State()
	assert.Equal(t,
		RepositoryState{
			State:           StateNone,
			Step:            "",
			Total:           "",
			HeadDescription: "(v1.1.0)",
			IsDetached:      true,
		},
		state,
	)
}

func TestStateOnPackedRefs(t *testing.T) {
	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("0123456789abcdef0123456789abcdef01234567\n"),
		},
		".git/packed-refs": &fstest.MapFile{
			Data: []byte("7c088a39dcd2dcda89f4dee1fd3eb41c1d34ea2f refs/heads/master\n0123456789abcdef0123456789abcdef01234567 refs/tags/v1.0.0\n"),
		},
	}

	git := testGitUtils("/Users/oriana/dev/kitsch", files)

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
	git := testGitUtils("", nil)

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
