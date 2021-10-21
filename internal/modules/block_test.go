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
		Join: " ",
		Modules: []ModuleSpec{
			{ID: "", Module: usernameMod},
			{ID: "", Module: promptMod},
		},
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
		Modules: []ModuleSpec{
			{ID: "", Module: usernameMod},
			{ID: "", Module: promptMod},
		},
		Join: " {{.PrevColors.FG}}{{.NextColors.FG}} ",
	}

	result := blockMod.Execute(testContext("jwalton"))
	assert.Equal(t, "jwalton redblue $ ", result.Text)
}

// TestBlockSubIDs verifies that the results of child modules can be indexed by ID.
func TestBlockSubIDs(t *testing.T) {
	usernameMod := UsernameModule{
		ShowAlways: true,
	}
	promptMod := PromptModule{}
	blockMod := BlockModule{
		CommonConfig: CommonConfig{
			Template: "{{ with .Data.Modules.a }}{{ .Data.Username }}{{ end }}",
		},
		Modules: []ModuleSpec{
			{ID: "a", Module: usernameMod},
			{ID: "", Module: promptMod},
		},
	}

	result := blockMod.Execute(testContext("oriana"))

	assert.Equal(t, "oriana", result.Text)
}
