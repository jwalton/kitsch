package gitutils

import (
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
)

func BenchmarkGetUpstream(b *testing.B) {
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

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		git.GetUpstream("master")
	}
}

func BenchmarkGetTagNameForHashAnnotatedTag(b *testing.B) {
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

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = git.GetTagNameForHash("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 100)
	}
}
