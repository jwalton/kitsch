package modules

import (
	"gopkg.in/yaml.v3"
)

// TextModule prints some text.
//
type TextModule struct {
	CommonConfig `yaml:",inline"`
	// Text is the text to print.
	Text string `yaml:"text"`
}

// Execute the module.
func (mod TextModule) Execute(context *Context) ModuleResult {
	data := map[string]interface{}{
		"Text": mod.Text,
	}

	return executeModule(context, mod.CommonConfig, data, mod.Style, mod.Text)
}

func init() {
	registerFactory("text", func(node *yaml.Node) (Module, error) {
		var module TextModule
		err := node.Decode(&module)
		return &module, err
	})
}
