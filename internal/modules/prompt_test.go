package modules

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/style"
	"github.com/stretchr/testify/assert"
)

func TestPrompt(t *testing.T) {
	mod := PromptModule{}
	env := &env.DummyEnv{}

	result := mod.Execute(env)
	assert.Equal(t, ModuleResult{
		Text: "$",
		Data: map[string]interface{}{
			"isRoot": false,
			"status": 0,
		},
		StartStyle: style.Style{},
		EndStyle:   style.Style{},
	}, result)
}

func TestRootPrompt(t *testing.T) {
	mod := PromptModule{}
	env := &env.DummyEnv{Root: true}

	result := mod.Execute(env)
	assert.Equal(t, ModuleResult{
		Text: "#",
		Data: map[string]interface{}{
			"isRoot": true,
			"status": 0,
		},
		StartStyle: style.Style{},
		EndStyle:   style.Style{},
	}, result)
}

func TestStyle(t *testing.T) {
	mod := PromptModule{
		CommonConfig: CommonConfig{
			Style: style.Style{FG: "blue"},
		},
		RootStyle: style.Style{FG: "red"},
	}

	stdEnv := &env.DummyEnv{Root: false}
	result := mod.Execute(stdEnv)

	rootEnv := &env.DummyEnv{Root: true}
	rootResult := mod.Execute(rootEnv)

	assert.Equal(t, "blue", result.StartStyle.FG)
	assert.Equal(t, "blue", result.EndStyle.FG)

	assert.Equal(t, "red", rootResult.StartStyle.FG)
	assert.Equal(t, "red", rootResult.EndStyle.FG)
}
