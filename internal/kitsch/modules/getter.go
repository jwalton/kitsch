package modules

import (
	"strings"

	"github.com/jwalton/kitsch-prompt/internal/kitsch/getters"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/log"
	"gopkg.in/yaml.v3"
)

// GetterModule executes a custom getter and returns the result.
//
// The `.Data` value returned from a custom module depends on the `As` configuration.
// If `As="text"`, then `.Data` will be a `{ Text: [string] }` object, containing
// the retrieved text (with leading and trailing whitespace automatically stripped).
// If `as` is any other value, then the `.Data` object will be the parsed result
// of the output.  For example if `as="json"`, and the returned value was
// '{"foo": "bar"}', then `.Data.foo` would be "bar".
//
type GetterModule struct {
	CommonConfig `yaml:",inline"`
	getterType   getters.GetterType

	// From is the source to get data from.  The meaning of "From" is based on
	// the provided "Type".
	From string `yaml:"from"`
	// As will determine how to interpret the result of the command.  One of
	// "text", "json", "toml", or "yaml".
	As getters.AsType `yaml:"as"`
	// Regex is a regular expression used to parse values out of the result of
	// the getter.  If specified, then "As" will be ignored.
	Regex string `yaml:"regex"`
	// Cache settings for the module.
	Cache getters.CacheSettings `yaml:"cache"`
}

type getterModuleTextResult struct {
	// Text is the text retrieved from the module.
	Text string
}

// Execute the module.
func (mod GetterModule) Execute(context *Context) ModuleResult {
	getter := getters.CustomGetter{
		Type:  mod.getterType,
		From:  mod.From,
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
		value = getterModuleTextResult{Text: text}
	} else {
		text = ""
	}

	return executeModule(context, mod.CommonConfig, value, mod.Style, text)
}

func init() {
	registerFactory("custom", func(node *yaml.Node) (Module, error) {
		module := GetterModule{getterType: getters.TypeCustom}
		err := node.Decode(&module)
		return &module, err
	})

	registerFactory("file", func(node *yaml.Node) (Module, error) {
		module := GetterModule{getterType: getters.TypeFile}
		err := node.Decode(&module)
		return &module, err
	})
}
