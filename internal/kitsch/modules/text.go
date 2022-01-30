package modules

import (
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas TextModule

// TextModule prints some text.
//
type TextModule struct {
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=text"`
	// Text is the text to print.
	Text string `yaml:"text" jsonschema:",required"`
}

type textModuleResult struct {
	// Text is the text to print. This is a mirror of the Text field in the configuration.
	Text string
}

// Execute the module.
func (mod TextModule) Execute(context *Context) ModuleResult {
	return ModuleResult{
		DefaultText: mod.Text,
		Data:        textModuleResult{Text: mod.Text},
	}
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
