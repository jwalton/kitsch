package modules

import (
	"strings"

	"github.com/jwalton/kitsch/internal/kitsch/getters"
	"github.com/jwalton/kitsch/internal/kitsch/log"
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas CustomModule

// CustomModule executes a custom getter and returns the result.
//
// The `.Data` value returned from a custom module depends on the `As` configuration.
// If `As="text"`, then `.Data` will be a `{ Text: [string] }` object, containing
// the retrieved text (with leading and trailing whitespace automatically stripped).
// If `as` is any other value, then the `.Data` object will be the parsed result
// of the output.  For example if `as="json"`, and the returned value was
// '{"foo": "bar"}', then `.Data.foo` would be "bar".
//
type CustomModule struct {
	// Type is the type of this module.
	Type getters.GetterType `yaml:"type" jsonschema:",required,enum=custom"`
	// Command is the command to run.
	Command string `yaml:"command" jsonschema:",required"`
	// As will determine how to interpret the result of the command.  One of
	// "text", "json", "toml", or "yaml".
	As getters.AsType `yaml:"as" jsonschema:",enum=text:json:toml:yaml"`
	// Regex is a regular expression used to parse values out of the result of
	// the getter.  If specified, then "As" will be ignored.
	Regex string `yaml:"regex"`
	// Cache settings for the module.
	Cache getters.CacheSettings `yaml:"cache" jsonschema:",ref"`
}

type customModuleTextResult struct {
	// Text is the text retrieved from the module.
	Text string
}

// Execute the module.
func (mod CustomModule) Execute(context *Context) ModuleResult {
	getter := getters.CustomGetter{
		Type:  mod.Type,
		From:  mod.Command,
		As:    mod.As,
		Regex: mod.Regex,
		Cache: mod.Cache,
	}

	value, err := getter.GetValue(context)
	if err != nil {
		log.Warn("Error executing custom module: ", err)
		value = ""
	}

	text, ok := value.(string)
	if ok {
		text = strings.TrimSpace(text)
		value = customModuleTextResult{Text: text}
	} else {
		text = ""
	}

	return ModuleResult{DefaultText: text, Data: value}
}

func init() {
	registerModule(
		"custom",
		registeredModule{
			jsonSchema: schemas.CustomModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := CustomModule{Type: getters.TypeCustom}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
