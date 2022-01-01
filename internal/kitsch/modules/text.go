package modules

import (
	"github.com/jwalton/kitsch-prompt/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas TextModule

// TextModule prints some text.
//
type TextModule struct {
	CommonConfig `yaml:",inline"`
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",enum=text"`
	// Text is the text to print.
	Text string `yaml:"text" jsonschema:",required"`
}

// Execute the module.
func (mod TextModule) Execute(context *Context) ModuleResult {
	data := map[string]interface{}{
		"Text": mod.Text,
	}

	return executeModule(context, mod.CommonConfig, data, mod.Style, mod.Text)
}

func init() {
	registerModule(
		"text",
		registeredModule{
			jsonSchema: schemas.TextModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := TextModule{Type: "text"}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
