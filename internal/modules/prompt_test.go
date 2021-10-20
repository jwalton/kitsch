package modules

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/styling"
	"github.com/stretchr/testify/assert"
)

func TestPrompt(t *testing.T) {
	mod := PromptModule{}
	context := testContext("jwalton")

	result := mod.Execute(context)

	assert.Equal(t, ModuleResult{
		Default: "$ ",
		Text:    "$ ",
		Data: map[string]interface{}{
			"IsRoot": false,
		},
		StartStyle: styling.CharacterColors{},
		EndStyle:   styling.CharacterColors{},
	}, result)
}

func TestRootPrompt(t *testing.T) {
	mod := PromptModule{}
	context := testContext("jwalton")
	context.Environment = &env.DummyEnv{Root: true}

	result := mod.Execute(context)

	assert.Equal(t, ModuleResult{
		Default: "# ",
		Text:    "# ",
		Data: map[string]interface{}{
			"IsRoot": true,
		},
		StartStyle: styling.CharacterColors{},
		EndStyle:   styling.CharacterColors{},
	}, result)
}

func TestStyle(t *testing.T) {
	mod := PromptModule{
		CommonConfig: CommonConfig{
			Style: "blue",
		},
		RootStyle: "red",
	}

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
