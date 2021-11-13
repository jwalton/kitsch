package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	blockMod := BlockModule{
		Join: " ",
		Modules: []ModuleSpec{
			{ID: "", Module: TextModule{Text: "hello"}},
			{ID: "", Module: TextModule{Text: "world"}},
		},
	}

	result := blockMod.Execute(testContext("jwalton"))
	assert.Equal(t, "hello world", result.Text)
}

func TestBlockStyles(t *testing.T) {
	blockMod := BlockModule{
		Modules: []ModuleSpec{
			{ID: "", Module: TextModule{
				CommonConfig: CommonConfig{Style: "red"},
				Text:         "hello",
			}},
			{ID: "", Module: TextModule{
				CommonConfig: CommonConfig{Style: "blue"},
				Text:         "world",
			}},
		},
		Join: " {{.PrevColors.FG}}{{.NextColors.FG}} ",
	}

	result := blockMod.Execute(testContext("jwalton"))
	assert.Equal(t, "hello redblue world", result.Text)
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
