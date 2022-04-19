package modules

import (
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch/internal/fileutils"
	"github.com/jwalton/kitsch/internal/kitsch/styling"
	"github.com/stretchr/testify/assert"
)

func TestModuleWrapperExecute(t *testing.T) {
	module := moduleWrapperFromYAML(heredoc.Doc(`
		type: text
		text: "test"
	`))

	result := module.Execute(newTestContext("jwalton"))

	assert.Equal(t,
		ModuleWrapperResult{
			Text:       "test",
			Data:       textModuleResult{Text: "test"},
			StartStyle: styling.CharacterColors{},
			EndStyle:   styling.CharacterColors{},
			Duration:   result.Duration,
		},
		result,
	)
}

func TestExecuteModuleWrapperWithTemplate(t *testing.T) {

	module := moduleWrapperFromYAML(heredoc.Doc(`
		type: text
		text: "Text Text"
		template: "--{{.Data.Text}}--"
	`))

	result := module.Execute(newTestContext("jwalton"))

	assert.Equal(t,
		ModuleWrapperResult{
			Text:       "--Text Text--",
			Data:       textModuleResult{Text: "Text Text"},
			StartStyle: styling.CharacterColors{},
			EndStyle:   styling.CharacterColors{},
			Duration:   result.Duration,
		},
		result,
	)
}

func TestExecuteModuleWithConditions(t *testing.T) {
	mod := moduleWrapperFromYAML(heredoc.Doc(`
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
