package modules

import (
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/stretchr/testify/assert"
)

func TestFileModule(t *testing.T) {
	context := newTestContext("jwalton")
	fsys := fstest.MapFS{
		"foo.json": {Data: []byte(`{"foo": "bar"}`)},
	}
	context.Directory = fileutils.NewDirectoryTestFS(context.Globals.CWD, fsys)

	mod := moduleFromYAMLMust(heredoc.Doc(`
		type: file
		from: foo.json
		as: json
		template: '{{.Data.foo}}'
	`))

	result := mod.Execute(context)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, result.Data)
	assert.Equal(t, "bar", result.Text)
}

func TestFileModuleRegex(t *testing.T) {
	context := newTestContext("jwalton")
	fsys := fstest.MapFS{
		"docker.txt": {Data: []byte(`Docker version 20.10.8, build 3967b7d`)},
	}
	context.Directory = fileutils.NewDirectoryTestFS(context.Globals.CWD, fsys)

	mod := moduleFromYAMLMust(heredoc.Doc(`
		type: file
		from: docker.txt
		regex: "^Docker version (.*), build .*$"
	`))

	result := mod.Execute(context)
	assert.Equal(t, getterModuleTextResult{Text: "20.10.8"}, result.Data)
	assert.Equal(t, "20.10.8", result.Text)
}

func TestFileModuleMissingFile(t *testing.T) {
	context := newTestContext("jwalton")
	fsys := fstest.MapFS{}
	context.Directory = fileutils.NewDirectoryTestFS(context.Globals.CWD, fsys)

	mod := moduleFromYAMLMust(heredoc.Doc(`
		type: file
		from: foo.json
		as: json
		template: '{{.Data.foo}}'
	`))

	result := mod.Execute(context)
	assert.Equal(t, getterModuleTextResult{Text: ""}, result.Data)
	assert.Equal(t, "", result.Text)
}
