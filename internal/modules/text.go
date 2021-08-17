package modules

import (
	"github.com/jwalton/kitsch-prompt/internal/env"
	"gopkg.in/yaml.v3"
)

// TextModule prints some text.
//
type TextModule struct {
	CommonConfig `yaml:",inline"`
	// Text is the text to print.
	Text string `yaml:"text"`
}

// Execute the directory module.
func (mod TextModule) Execute(env env.Env) ModuleResult {
	data := map[string]interface{}{
		"text": mod.Text,
	}

	return executeModule(mod.CommonConfig, data, mod.Style, mod.Text)
}

func init() {
	registerFactory("text", func(node *yaml.Node) (Module, error) {
		var module TextModule
		err := node.Decode(&module)
		return &module, err
	})
}
