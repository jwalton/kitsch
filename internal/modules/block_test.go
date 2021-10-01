package modules

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/style"
	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	usernameMod := UsernameModule{
		ShowAlways: true,
	}
	promptMod := PromptModule{}
	blockMod := BlockModule{
		Modules: ModuleList{[]Module{usernameMod, promptMod}},
		Join:    " ",
	}

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
		},
	}

	result := blockMod.Execute(env)
	assert.Equal(t, "jwalton $ ", result.Text)
}

func TestBlockStyles(t *testing.T) {
	usernameMod := UsernameModule{
		CommonConfig: CommonConfig{Style: style.Style{FG: "red"}},
		ShowAlways:   true,
	}
	promptMod := PromptModule{
		CommonConfig: CommonConfig{Style: style.Style{FG: "blue"}},
	}

	blockMod := BlockModule{
		Modules: ModuleList{[]Module{usernameMod, promptMod}},
		Join:    " {{.prevStyle.FG}}{{.nextStyle.FG}} ",
	}

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
		},
	}

	result := blockMod.Execute(env)
	assert.Equal(t, "jwalton redblue $ ", result.Text)
}
