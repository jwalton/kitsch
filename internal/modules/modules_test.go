package modules

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/internal/style"
	"github.com/stretchr/testify/assert"
)

func TestExecuteModule(t *testing.T) {
	data := map[string]interface{}{}

	result := executeModule(
		testContext("jwalton"),
		CommonConfig{},
		data,
		"",
		"test",
	)

	assert.Equal(t,
		ModuleResult{
			Default:    "test",
			Text:       "test",
			Data:       data,
			StartStyle: style.CharacterColors{},
			EndStyle:   style.CharacterColors{},
		},
		result,
	)
}

func TestExecuteModuleWithTemplate(t *testing.T) {
	data := map[string]interface{}{
		"Text": "Text Text",
	}
	moduleStyle := ""

	result := executeModule(
		testContext("jwalton"),
		CommonConfig{
			Template: "--{{.Data.Text}}--",
		},
		data,
		moduleStyle,
		"test",
	)

	assert.Equal(t,
		ModuleResult{
			Default:    "test",
			Text:       "--Text Text--",
			Data:       data,
			StartStyle: style.CharacterColors{},
			EndStyle:   style.CharacterColors{},
		},
		result,
	)
}
