package modules

import (
	"testing"
	"testing/fstest"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/condition"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/getters"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/projects"
	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	context := testContext("jwalton")

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
			Conditions:            condition.Conditions{IfExtensions: []string{"txt"}},
			ToolSymbol:            "txt",
			PackageManagerSymbol:  "txt",
			ToolVersion:           getters.CustomGetter{Type: "file", From: "tool.txt"},
			PackageManagerVersion: getters.CustomGetter{Type: "file", From: "pm.txt"},
			PackageVersion:        getters.CustomGetter{Type: "file", From: "package.txt"},
		},
	}

	mod := ProjectModule{
		CommonConfig: CommonConfig{
			Template: "{{ .Data.ToolVersion }} / {{ .Data.PackageManagerVersion }} / {{ .Data.PackageVersion }}",
		},
	}

	result := mod.Execute(context)

	assert.Equal(t, "v1.0.0 / v2.0.0 / v3.0.0", result.Text)
}
