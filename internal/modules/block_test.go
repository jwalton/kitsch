package modules

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/style"
	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	usernameMod := NewUsernameModule(UsernameConfig{
		ShowAlways: true,
	})
	promptMod := NewPromptModule(PromptConfig{})
	blockMod := NewBlockModule(BlockConfig{
		Modules: []Module{usernameMod, promptMod},
		Join:    " ",
	})

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
		},
	}

	result := blockMod.Execute(env)
	assert.Equal(t, "jwalton $", result.Text)
}

func TestBlockStyles(t *testing.T) {
	usernameMod := NewUsernameModule(UsernameConfig{
		CommonConfig: CommonConfig{Style: style.Style{FG: "red"}},
		ShowAlways:   true,
	})
	promptMod := NewPromptModule(PromptConfig{
		CommonConfig: CommonConfig{Style: style.Style{FG: "blue"}},
	})

	blockMod := NewBlockModule(BlockConfig{
		Modules: []Module{usernameMod, promptMod},
		Join:    " {{.prevStyle.FG}}{{.nextStyle.FG}} ",
	})

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
		},
	}

	result := blockMod.Execute(env)
	assert.Equal(t, "jwalton redblue $", result.Text)
}
