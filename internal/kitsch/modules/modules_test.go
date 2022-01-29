package modules

import (
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch/internal/fileutils"
	"github.com/jwalton/kitsch/internal/kitsch/styling"
	"github.com/stretchr/testify/assert"
)

func TestExecuteModule(t *testing.T) {
	data := map[string]interface{}{}

	result := executeModule(
		newTestContext("jwalton"),
		CommonConfig{},
		data,
		"",
		"test",
	)

	assert.Equal(t,
		ModuleResult{
			Text:       "test",
			Data:       data,
			StartStyle: styling.CharacterColors{},
			EndStyle:   styling.CharacterColors{},
		},
		result,
	)
}

func TestExecuteModuleWithTemplate(t *testing.T) {
	data := map[string]interface{}{
		"Text": "Text Text",
	}
	moduleStyle := ""

	result := executeModule(
		newTestContext("jwalton"),
		CommonConfig{
			Template: "--{{.Data.Text}}--",
		},
		data,
		moduleStyle,
		"test",
	)

	assert.Equal(t,
		ModuleResult{
			Text:       "--Text Text--",
			Data:       data,
			StartStyle: styling.CharacterColors{},
			EndStyle:   styling.CharacterColors{},
		},
		result,
	)
}

func TestExecuteModuleWithConditions(t *testing.T) {
	mod := moduleSpecFromYAML(heredoc.Doc(`
		type: text
		conditions:
		  ifFiles:  ['helm']
		text: Hello World
	`))

	context := newTestContext("jwalton")

	context.Directory = fileutils.NewDirectoryTestFS("/foo/bar", fstest.MapFS{
		"helm": &fstest.MapFile{
			Data: []byte("blahblahblah"),
		},
	})
	result := mod.Execute(context)
	assert.Equal(t, "Hello World", result.Text)

	context.Directory = fileutils.NewDirectoryTestFS("/foo/bar", fstest.MapFS{})
	result2 := mod.Execute(context)
	assert.Equal(t, "", result2.Text)
}
