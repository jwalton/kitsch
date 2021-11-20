package modules

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/styling"
	"github.com/stretchr/testify/assert"
)

func TestPrompt(t *testing.T) {
	mod := moduleFromYAMLMust("{type: prompt}")

	context := newTestContext("jwalton")

	result := mod.Execute(context)

	assert.Equal(t, ModuleResult{
		Text: "$ ",
		Data: promptModuleData{
			PromptString: "$ ",
			PromptStyle:  "",
			ViCmdMode:    false,
		},
		StartStyle: styling.CharacterColors{},
		EndStyle:   styling.CharacterColors{},
	}, result)
}

func TestRootPrompt(t *testing.T) {
	mod := moduleFromYAMLMust("{type: prompt}")

	context := newTestContext("jwalton")
	context.Globals.IsRoot = true

	result := mod.Execute(context)

	assert.Equal(t, ModuleResult{
		Text: "# ",
		Data: promptModuleData{
			PromptString: "# ",
			PromptStyle:  "",
			ViCmdMode:    false,
		},
		StartStyle: styling.CharacterColors{},
		EndStyle:   styling.CharacterColors{},
	}, result)
}

func TestStyle(t *testing.T) {
	mod := moduleFromYAMLMust(heredoc.Doc(`
		type: prompt
		style: blue
		rootStyle: red
	`))

	context := newTestContext("jwalton")

	result := mod.Execute(context)

	context.Globals.IsRoot = true
	rootResult := mod.Execute(context)

	assert.Equal(t, "blue", result.StartStyle.FG)
	assert.Equal(t, "blue", result.EndStyle.FG)

	assert.Equal(t, "red", rootResult.StartStyle.FG)
	assert.Equal(t, "red", rootResult.EndStyle.FG)
}
