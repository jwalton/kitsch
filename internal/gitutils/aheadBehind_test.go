package gitutils

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/jwalton/kitsch/internal/billyutils"
	"github.com/stretchr/testify/assert"
)

// testGitUtils creates a new gitUtils for unit testing.
func testGitUtils(repoRoot string, fsys fs.FS) *gitUtils {
	var storer *filesystem.Storage
	if fsys != nil {
		repositoryFs, err := billyutils.FsToBilly(fsys)
		if err != nil {
			panic(err)
		}
		dotGitFs, err := repositoryFs.Chroot(".git")
		if err != nil {
			panic(err)
		}

		storer = filesystem.NewStorage(dotGitFs, cache.NewObjectLRUDefault())
	}

	return &gitUtils{
		pathToGit: "git",
		fsys:      fsys,
		repoRoot:  repoRoot,
		storer:    storer,
	}
}

func TestGetUpstream(t *testing.T) {
	config := heredoc.Doc(`
		[branch "master"]
			remote = origin
			merge = refs/heads/master
		[branch "feature/projects"]
			remote = spooky
			merge = refs/heads/feature/oldprojects
	`)

	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("ref: refs/heads/master\n"),
		},
		".git/config": &fstest.MapFile{
			Data: []byte(config),
		},
	}

	git := testGitUtils("/Users/oriana/dev/kitsch", files)

	assert.Equal(t,
		"origin/master",
		git.GetUpstream("master"),
	)

	assert.Equal(t,
		"spooky/feature/oldprojects",
		git.GetUpstream("feature/projects"),
	)

	assert.Equal(t,
		"",
		git.GetUpstream("banana"),
	)
}

func TestGetUpstreamNoConfig(t *testing.T) {
	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("ref: refs/heads/master\n"),
		},
	}

	git := testGitUtils("/Users/oriana/dev/kitsch", files)

	assert.Equal(t,
		"",
		git.GetUpstream("feature/projects"),
	)
}
