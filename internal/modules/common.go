package modules

import (
	"bytes"
	"text/template"

	styleLib "github.com/jwalton/kitsch-prompt/internal/style"
)

// CommonConfig is common configuration for all modules.
type CommonConfig struct {
	Style    styleLib.Style `yaml:"style"`
	Template string         `yaml:"template"`
	// // Prefix is a string to add to the start of the output of this module.
	// Prefix   string `yaml:"prefix"`
	// // Style to apply to the prefix.
	// PrefixStyle styleLib.Style `yaml:"prefixStyle"`
	// // Suffix is a string to add to the end of the output of this module.
	// Suffix   string `yaml:"suffix"`
	// // Style to apply to the suffix.
	// SuffixStyle styleLib.Style `yaml:"suffixStyle"`
}

// func readCommonConfig(data map[string]interface{}, config *CommonConfig) error {
// 	var err error = nil

// 	template, ok := data["template"].(string)
// 	if ok {
// 		config.Template = template
// 	}

// 	rawStyle, ok := data["style"]
// 	if ok {
// 		err = config.Style.UnmarshalInterface(rawStyle)
// 	}

// 	return err
// }

// executeModule is called to execute a module.  This handles "common" stuff that
// all modules do, like calling templates.
func executeModule(
	config CommonConfig,
	data map[string]interface{},
	style styleLib.Style,
	defaultText string,
) ModuleResult {
	text := defaultText
	var startStyle styleLib.Style
	var endStyle styleLib.Style

	// prefix, prefixStartStyle, _, prefixErr := config.PrefixStyle.Apply(config.Prefix)
	// suffix, _, suffixEndStyle, _, suffixErr := config.SuffixStyle.Apply(config.Suffix)

	if config.Template != "" {
		tmpl, err := template.New("module-template").Parse(config.Template)
		// FIX: Should add this error to a list of warnings for this module.
		if err == nil {
			text = templateToString(tmpl, data)
		}
	}

	if !style.IsEmpty() {
		// TODO: Do something useful with this error?
		text, startStyle, endStyle, _ = style.Apply(text)
	}

	return ModuleResult{
		Text:       text,
		Data:       data,
		StartStyle: startStyle,
		EndStyle:   endStyle,
	}
}

func templateToString(template *template.Template, data interface{}) string {
	var b bytes.Buffer
	err := template.Execute(&b, data)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

// defaultString returns value if it is non-empty, or def otherwise.
func defaultString(value string, def string) string {
	if value != "" {
		return value
	}
	return def
}

// defaultNumber returns the passed in value if it is non-0, or def otherwise.
func defaultNumber(value int, def int) int {
	if value != 0 {
		return value
	}
	return def
}
