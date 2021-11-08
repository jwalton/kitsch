package condition

import (
	"testing"
	"testing/fstest"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/stretchr/testify/assert"
)

func TestIfFiles(t *testing.T) {
	directory := fileutils.NewDirectoryTestFS(
		"/foo/bar",
		fstest.MapFS{
			"version.txt": &fstest.MapFile{
				Data: []byte("v1.0.0\n"),
			},
			"file.go": &fstest.MapFile{
				Data: []byte("package main\n"),
			},
		},
	)

	condition := Condition{IfFiles: []string{"version.txt"}}
	assert.Equal(t, true, condition.Matches(directory))

	condition = Condition{IfFiles: []string{"nothere.txt"}}
	assert.Equal(t, false, condition.Matches(directory))

	condition = Condition{IfExtensions: []string{"go"}}
	assert.Equal(t, true, condition.Matches(directory))

	condition = Condition{IfExtensions: []string{".go"}}
	assert.Equal(t, true, condition.Matches(directory))

	condition = Condition{IfExtensions: []string{".js"}}
	assert.Equal(t, false, condition.Matches(directory))

	condition = Condition{IfFiles: []string{"nothere.txt"}, IfExtensions: []string{".js"}}
	assert.Equal(t, false, condition.Matches(directory))

	condition = Condition{IfFiles: []string{"version.txt"}, IfExtensions: []string{".js"}}
	assert.Equal(t, true, condition.Matches(directory))

	condition = Condition{IfFiles: []string{"nothere.txt"}, IfExtensions: []string{".go"}}
	assert.Equal(t, true, condition.Matches(directory))

	condition = Condition{IfFiles: []string{"version.txt"}, IfExtensions: []string{".go"}}
	assert.Equal(t, true, condition.Matches(directory))
}
