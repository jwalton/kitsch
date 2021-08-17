package modules

import (
	"github.com/jwalton/kitsch-prompt/internal/env"
	"gopkg.in/yaml.v3"
)

// StringModule shows a literal string.
//
type StringModule struct {
	CommonConfig `yaml:",inline"`
	// Value is the value to show.
	Value string
}

// Execute the prompt module.
func (mod StringModule) Execute(env env.Env) ModuleResult {
	data := map[string]interface{}{
		"value": mod.Value,
	}

	return executeModule(mod.CommonConfig, data, mod.Style, mod.Value)
}

func init() {
	registerFactory("string", func(node *yaml.Node) (Module, error) {
		var module StringModule
		err := node.Decode(&module)
		return &module, err
	})
}
