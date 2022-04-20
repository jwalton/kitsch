package modules

import (
	"strings"

	"github.com/jwalton/kitsch/internal/kitsch/getters"
	"github.com/jwalton/kitsch/internal/kitsch/log"
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas FileModule

// FileModule executes a custom getter and returns the result.
//
// The `.Data` value returned from a custom module depends on the `As` configuration.
// If `As="text"`, then `.Data` will be a `{ Text: [string] }` object, containing
// the retrieved text (with leading and trailing whitespace automatically stripped).
// If `as` is any other value, then the `.Data` object will be the parsed result
// of the output.  For example if `as="json"`, and the returned value was
// '{"foo": "bar"}', then `.Data.foo` would be "bar".
//
type FileModule struct {
	// Type is the type of this module.
	Type getters.GetterType `yaml:"type" jsonschema:",required,enum=file"`
	// From is the file to read data from.
	File string `yaml:"file" jsonschema:",required"`
	// As will determine how to interpret the contents of the file.  One of
	// "text", "json", "toml", or "yaml".
	As getters.AsType `yaml:"as" jsonschema:",enum=text:json:toml:yaml"`
	// Regex is a regular expression used to parse values out of the result of
	// the getter.  If specified, then "As" will be ignored.
	Regex string `yaml:"regex"`
}

type fileModuleTextResult struct {
	// Text is the text retrieved from the module.
	Text string
}

// Execute the module.
func (mod FileModule) Execute(context *Context) ModuleResult {
	getter := getters.CustomGetter{
		Type:  mod.Type,
		From:  mod.File,
		As:    mod.As,
		Regex: mod.Regex,
	}

	value, err := getter.GetValue(context)
	if err != nil {
		log.Warn("Error executing file module: ", err)
		value = ""
	}

	text, ok := value.(string)
	if ok {
		text = strings.TrimSpace(text)
		value = fileModuleTextResult{Text: text}
	} else {
		text = ""
	}

	return ModuleResult{DefaultText: text, Data: value}
}

func init() {
	registerModule(
		"file",
		registeredModule{
			jsonSchema: schemas.FileModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := FileModule{Type: getters.TypeFile}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
