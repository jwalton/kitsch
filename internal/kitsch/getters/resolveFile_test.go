package getters

import (
	"path/filepath"
	"strings"
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
	expected := strings.Replace("/users/jwalton/foo/bar/baz", "/", string(filepath.Separator), -1)

	assert.Equal(t, expected, result)
}
