package modules

import (
	"testing"

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

	result := blockMod.Execute(testContext("jwalton"))
	assert.Equal(t, "jwalton $ ", result.Text)
}

func TestBlockStyles(t *testing.T) {
	usernameMod := UsernameModule{
		CommonConfig: CommonConfig{Style: "red"},
		ShowAlways:   true,
	}
	promptMod := PromptModule{
		CommonConfig: CommonConfig{Style: "blue"},
	}

	blockMod := BlockModule{
		Modules: ModuleList{[]Module{usernameMod, promptMod}},
		Join:    " {{.PrevColors.FG}}{{.NextColors.FG}} ",
	}

	result := blockMod.Execute(testContext("jwalton"))
	assert.Equal(t, "jwalton redblue $ ", result.Text)
}
