package condition

import (
	"testing"
	"testing/fstest"

	"github.com/jwalton/kitsch/internal/fileutils"
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

	conditions := Conditions{IfFiles: []string{"version.txt"}}
	assert.Equal(t, true, conditions.Matches(directory))

	conditions = Conditions{IfFiles: []string{"nothere.txt"}}
	assert.Equal(t, false, conditions.Matches(directory))

	conditions = Conditions{IfExtensions: []string{"go"}}
	assert.Equal(t, true, conditions.Matches(directory))

	conditions = Conditions{IfExtensions: []string{".go"}}
	assert.Equal(t, true, conditions.Matches(directory))

	conditions = Conditions{IfExtensions: []string{".js"}}
	assert.Equal(t, false, conditions.Matches(directory))

	conditions = Conditions{IfFiles: []string{"nothere.txt"}, IfExtensions: []string{".js"}}
	assert.Equal(t, false, conditions.Matches(directory))

	conditions = Conditions{IfFiles: []string{"version.txt"}, IfExtensions: []string{".js"}}
	assert.Equal(t, true, conditions.Matches(directory))

	conditions = Conditions{IfFiles: []string{"nothere.txt"}, IfExtensions: []string{".go"}}
	assert.Equal(t, true, conditions.Matches(directory))

	conditions = Conditions{IfFiles: []string{"version.txt"}, IfExtensions: []string{".go"}}
	assert.Equal(t, true, conditions.Matches(directory))
}
