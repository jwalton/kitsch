package gitutils

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestHeadOnBranch(t *testing.T) {
	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("ref: refs/heads/master\n"),
		},
		".git/refs/heads/master": &fstest.MapFile{
			Data: []byte("7c088a39dcd2dcda89f4dee1fd3eb41c1d34ea2f\n"),
		},
	}

	git := testGitUtils("/Users/oriana/dev/kitsch", files)

	state, err := git.Head(100)
	assert.Nil(t, err)
	assert.Equal(t,
		HeadInfo{
			Description: "master",
			Detached:    false,
			Hash:        "7c088a39dcd2dcda89f4dee1fd3eb41c1d34ea2f",
			IsTag:       false,
		},
		state,
	)
}

func TestHeadDetached(t *testing.T) {
	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("0123456789abcdef0123456789abcdef01234567\n"),
		},
		".git/refs/heads/master": &fstest.MapFile{
			Data: []byte("7c088a39dcd2dcda89f4dee1fd3eb41c1d34ea2f\n"),
		},
	}

	git := testGitUtils("/Users/oriana/dev/kitsch", files)

	state, err := git.Head(100)
	assert.Nil(t, err)
	assert.Equal(t,
		HeadInfo{
			Description: "(0123456...)",
			Detached:    true,
			Hash:        "0123456789abcdef0123456789abcdef01234567",
			IsTag:       false,
		},
		state,
	)
}

func TestHeadOnTag(t *testing.T) {
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

	state, err := git.Head(100)
	assert.Nil(t, err)
	assert.Equal(t,
		HeadInfo{
			Description: "(v1.0.0)",
			Detached:    true,
			Hash:        "0123456789abcdef0123456789abcdef01234567",
			IsTag:       true,
		},
		state,
	)

	// Disable searching for tags.
	state, err = git.Head(0)
	assert.Nil(t, err)
	assert.Equal(t,
		HeadInfo{
			Description: "(0123456...)",
			Detached:    true,
			Hash:        "0123456789abcdef0123456789abcdef01234567",
			IsTag:       false,
		},
		state,
	)

}

func TestHeadOnAnnotatedTag(t *testing.T) {
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

	state, err := git.Head(100)
	assert.Nil(t, err)
	assert.Equal(t,
		HeadInfo{
			Description: "(v1.1.0)",
			Detached:    true,
			Hash:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			IsTag:       true,
		},
		state,
	)
}

func TestHeadOnPackedRefs(t *testing.T) {
	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("0123456789abcdef0123456789abcdef01234567\n"),
		},
		".git/packed-refs": &fstest.MapFile{
			Data: []byte("7c088a39dcd2dcda89f4dee1fd3eb41c1d34ea2f refs/heads/master\n0123456789abcdef0123456789abcdef01234567 refs/tags/v1.0.0\n"),
		},
	}

	git := testGitUtils("/Users/oriana/dev/kitsch", files)

	state, err := git.Head(100)
	assert.Nil(t, err)
	assert.Equal(t,
		HeadInfo{
			Description: "(v1.0.0)",
			Detached:    true,
			Hash:        "0123456789abcdef0123456789abcdef01234567",
			IsTag:       true,
		},
		state,
	)
}

func TestHeadInNonGitRepo(t *testing.T) {
	git := testGitUtils("", nil)
	_, err := git.Head(100)
	assert.EqualError(t, err, "no git repo found")
}
