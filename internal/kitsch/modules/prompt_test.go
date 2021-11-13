package modules

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/env"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/styling"
	"github.com/stretchr/testify/assert"
)

func TestPrompt(t *testing.T) {
	mod, err := moduleFromYAML("{type: prompt}")
	assert.NoError(t, err)

	context := testContext("jwalton")

	result := mod.Execute(context)

	assert.Equal(t, ModuleResult{
		Text: "$ ",
		Data: map[string]interface{}{
			"IsRoot": false,
		},
		StartStyle: styling.CharacterColors{},
		EndStyle:   styling.CharacterColors{},
	}, result)
}

func TestRootPrompt(t *testing.T) {
	mod, err := moduleFromYAML("{type: prompt}")
	assert.NoError(t, err)

	context := testContext("jwalton")
	context.Environment = &env.DummyEnv{Root: true}

	result := mod.Execute(context)

	assert.Equal(t, ModuleResult{
		Text: "# ",
		Data: map[string]interface{}{
			"IsRoot": true,
		},
		StartStyle: styling.CharacterColors{},
		EndStyle:   styling.CharacterColors{},
	}, result)
}

func TestStyle(t *testing.T) {
	mod, err := moduleFromYAML(heredoc.Doc(`
		type: prompt
		style: blue
		rootStyle: red
	`))
	assert.NoError(t, err)

	context := testContext("jwalton")
	context.Environment = &env.DummyEnv{Root: false}

	result := mod.Execute(context)

	context.Environment = &env.DummyEnv{Root: true}
	rootResult := mod.Execute(context)

	assert.Equal(t, "blue", result.StartStyle.FG)
	assert.Equal(t, "blue", result.EndStyle.FG)

	assert.Equal(t, "red", rootResult.StartStyle.FG)
	assert.Equal(t, "red", rootResult.EndStyle.FG)
}
