package gitutils

import (
	"bytes"
	"compress/zlib"
	"strconv"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestGetStashCount(t *testing.T) {
	stashes := `0000000000000000000000000000000000000000 b0592fb675bd471541928aa8c9900ba76f748ac8 Jason Walton <dev@lucid.thedreaming.org> 1634679248 -0400	WIP on master: d5f7bc9 feat: Rewrite internal/style. Split template data into Data and Global. Add support for custom colors.
b0592fb675bd471541928aa8c9900ba76f748ac8 9f139fbbde7200508adecc1b9adad67e99204ae3 Jason Walton <dev@lucid.thedreaming.org> 1634734842 -0400	WIP on master: 354da14 fix(style): Fix style template function to work with non-strings.
9f139fbbde7200508adecc1b9adad67e99204ae3 b5327c73a4c87d535b88b2238c18b27ffe091eab Jason Walton <dev@lucid.thedreaming.org> 1636838284 -0500	WIP on master: 6e0e0a8 feat: Add vi mode support.
`

	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("ref: refs/heads/master\n"),
		},
		".git/logs/refs/stash": &fstest.MapFile{
			Data: []byte(stashes),
		},
	}

	git := &GitUtils{
		pathToGit: "git",
		fsys:      files,
		RepoRoot:  "/Users/oriana/dev/kitsch-prompt",
	}

	result, err := git.GetStashCount()
	assert.Nil(t, err)
	assert.Equal(t, 3, result)
}

func TestGetStashCountNoStashes(t *testing.T) {
	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("ref: refs/heads/master\n"),
		},
	}

	git := &GitUtils{
		pathToGit: "git",
		fsys:      files,
		RepoRoot:  "/Users/oriana/dev/kitsch-prompt",
	}

	result, err := git.GetStashCount()
	assert.Nil(t, err)
	assert.Equal(t, 0, result)
}

// generateGitObject generates a git object with the given content.
func generateGitObject(objectType string, content string) []byte {
	contentBytes := []byte(content)
	out := bytes.Buffer{}
	out.WriteString(objectType + " ")
	out.WriteString(strconv.Itoa(len(contentBytes)))
	out.WriteByte(0)
	out.Write(contentBytes)

	compressedContent := bytes.Buffer{}
	w := zlib.NewWriter(&compressedContent)
	_, _ = w.Write(out.Bytes())
	w.Close()

	return compressedContent.Bytes()
}

func TestReadObject(t *testing.T) {
	files := fstest.MapFS{
		".git/objects/b0/592fb675bd471541928aa8c9900ba76f748ac8": &fstest.MapFile{
			Data: generateGitObject("blob", "hello world"),
		},
	}

	git := &GitUtils{
		pathToGit: "git",
		fsys:      files,
		RepoRoot:  "/Users/oriana/dev/kitsch-prompt",
	}

	objectType, content, err := git.ReadObject("b0592fb675bd471541928aa8c9900ba76f748ac8")
	assert.Nil(t, err)
	assert.Equal(t, "blob", objectType)
	assert.Equal(t, "hello world", string(content))

	content, err = git.ReadObjectOfType("blob", "b0592fb675bd471541928aa8c9900ba76f748ac8")
	assert.Nil(t, err)
	assert.Equal(t, "hello world", string(content))

	_, err = git.ReadObjectOfType("tag", "b0592fb675bd471541928aa8c9900ba76f748ac8")
	assert.Equal(t, ErrIncorrectType, err)
}
