package modules

import (
	"testing"
	"testing/fstest"
	"time"

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

type sleepModule struct {
	// Type is the type of this module.
	Type string
	// Duration is the time to sleep in milliseconds.
	Duration int64
	// Text is the output of this module.
	Text string
}

// Execute the module.
func (mod sleepModule) Execute(context *Context) ModuleResult {
	time.Sleep(time.Duration(mod.Duration) * time.Millisecond)

	return ModuleResult{
		DefaultText: mod.Text,
		Data:        nil,
	}
}

func TestExecuteModuleWithTimeout(t *testing.T) {
	mod := ModuleWrapper{
		config: CommonConfig{
			Type:    "sleep",
			Timeout: 10,
		},
		Module: sleepModule{
			Type:     "sleep",
			Duration: 1000,
			Text:     "Hello World",
		},
	}

	context := newTestContext("jwalton")
	result := mod.Execute(context)

	// Should have no output, because it should have timed out.
	assert.Equal(t, "", result.Text)
}
