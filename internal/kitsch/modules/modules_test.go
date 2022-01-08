package modules

import (
	"testing"

	"github.com/jwalton/kitsch/internal/kitsch/styling"
	"github.com/stretchr/testify/assert"
)

func TestExecuteModule(t *testing.T) {
	data := map[string]interface{}{}

	result := executeModule(
		newTestContext("jwalton"),
		CommonConfig{},
		data,
		"",
		"test",
	)

	assert.Equal(t,
		ModuleResult{
			Text:       "test",
			Data:       data,
			StartStyle: styling.CharacterColors{},
			EndStyle:   styling.CharacterColors{},
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
		newTestContext("jwalton"),
		CommonConfig{
			Template: "--{{.Data.Text}}--",
		},
		data,
		moduleStyle,
		"test",
	)

	assert.Equal(t,
		ModuleResult{
			Text:       "--Text Text--",
			Data:       data,
			StartStyle: styling.CharacterColors{},
			EndStyle:   styling.CharacterColors{},
		},
		result,
	)
}
