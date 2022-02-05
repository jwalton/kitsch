package modules

import (
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch/internal/fileutils"
	"github.com/jwalton/kitsch/internal/kitsch/condition"
	"github.com/jwalton/kitsch/internal/kitsch/getters"
	"github.com/jwalton/kitsch/internal/kitsch/projects"
	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	context := newTestContext("jwalton")

	fsys := fstest.MapFS{
		"tool.txt": &fstest.MapFile{
			Data: []byte("v1.0.0\n"),
		},
		"pm.txt": &fstest.MapFile{
			Data: []byte("v2.0.0\n"),
		},
		"package.txt": &fstest.MapFile{
			Data: []byte("v3.0.0\n"),
		},
	}
	context.Directory = fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	context.ProjectTypes = []projects.ProjectType{
		{
			Name:                  "txt",
			Conditions:            &condition.Conditions{IfExtensions: []string{"txt"}},
			ToolSymbol:            "txt",
			PackageManagerSymbol:  "txt",
			ToolVersion:           []getters.Getter{getters.CustomGetter{Type: getters.TypeFile, From: "tool.txt"}},
			PackageManagerVersion: []getters.Getter{getters.CustomGetter{Type: getters.TypeFile, From: "pm.txt"}},
			PackageVersion:        []getters.Getter{getters.CustomGetter{Type: getters.TypeFile, From: "package.txt"}},
		},
	}

	mod := moduleWrapperFromYAML(heredoc.Doc(`
		type: project
		template: "{{ .Data.ToolVersion }} / {{ .Data.PackageManagerVersion }} / {{ .Data.PackageVersion }}"
	`))

	result := mod.Execute(context)

	assert.Equal(t, "v1.0.0 / v2.0.0 / v3.0.0", result.Text)
}
