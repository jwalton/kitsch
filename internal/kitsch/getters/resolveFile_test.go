package getters

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestResolveFile(t *testing.T) {
	context := makeTestGetterContext(fstest.MapFS{})
	context.env = map[string]string{
		"VAR1": "foo",
		"VAR2": "bar",
	}

	result := resolveFile(context, "~/${VAR1}/$VAR2/baz")
	assert.Equal(t, "/users/jwalton/foo/bar/baz", result)
}
