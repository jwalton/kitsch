package modules

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/internal/style"
	"github.com/stretchr/testify/assert"
)

func TestExecuteModule(t *testing.T) {
	data := map[string]interface{}{}
	moduleStyle := style.Style{}

	result := executeModule(
		CommonConfig{},
		data,
		moduleStyle,
		"test",
	)

	assert.Equal(t, ModuleResult{
		Text:       "test",
		Data:       data,
		StartStyle: moduleStyle,
		EndStyle:   moduleStyle,
	}, result)
}

func TestExecuteModuleWithTemplate(t *testing.T) {
	data := map[string]interface{}{
		"text": "Text Text",
	}
	moduleStyle := style.Style{}

	result := executeModule(
		CommonConfig{
			Template: "--{{.text}}--",
		},
		data,
		moduleStyle,
		"test",
	)

	assert.Equal(t, ModuleResult{
		Text:       "--Text Text--",
		Data:       data,
		StartStyle: moduleStyle,
		EndStyle:   moduleStyle,
	}, result)
}
