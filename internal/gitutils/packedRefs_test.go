package gitutils

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

var packedRefsData = heredoc.Doc(`
	# This is a comment
	b0d69e9aae0d71548e41c08ccedd00e791bc827c refs/heads/master
	b0d69e9aae0d71548e41c08ccedd00e791bc827c refs/remotes/origin/master
	3a2259e5b6b172c88028b408bf0a22848dbfc9ef refs/tags/v0.0.1
	ed34073090b55039bcae69324ffd251dd1117b2b refs/tags/v0.0.2
`)

func TestForEachPackedRef(t *testing.T) {
	var refs []string

	forEachPackedRef([]byte(packedRefsData), func(hash string, ref string) bool {
		refs = append(refs, hash)
		refs = append(refs, ref)
		return true
	})

	assert.Equal(t, []string{
		"b0d69e9aae0d71548e41c08ccedd00e791bc827c", "refs/heads/master",
		"b0d69e9aae0d71548e41c08ccedd00e791bc827c", "refs/remotes/origin/master",
		"3a2259e5b6b172c88028b408bf0a22848dbfc9ef", "refs/tags/v0.0.1",
		"ed34073090b55039bcae69324ffd251dd1117b2b", "refs/tags/v0.0.2",
	}, refs)
}

func TestForEachPackedRefStop(t *testing.T) {
	var refs []string

	forEachPackedRef([]byte(packedRefsData), func(hash string, ref string) bool {
		refs = append(refs, hash)
		refs = append(refs, ref)
		return ref != "refs/tags/v0.0.1"
	})

	// Should only fetch refs until we get to refs/tags/v0.0.1.
	assert.Equal(t, []string{
		"b0d69e9aae0d71548e41c08ccedd00e791bc827c", "refs/heads/master",
		"b0d69e9aae0d71548e41c08ccedd00e791bc827c", "refs/remotes/origin/master",
		"3a2259e5b6b172c88028b408bf0a22848dbfc9ef", "refs/tags/v0.0.1",
	}, refs)
}
