package modules

import (
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas PromptModule

// PromptModule shows a prompt to the user.
//
// The prompt module displays a "$", or a "#" if the current user is root.
//
type PromptModule struct {
	CommonConfig `yaml:",inline"`
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",enum=prompt"`
	// Prompt is what to display as the prompt.  Defaults to "$ ".
	Prompt string `yaml:"prompt"`
	// RootPrompt is what to display as the prompt if the current user is root.  Defaults to "# ".
	RootPrompt string `yaml:"rootPrompt"`
	// RootStyle will be used in place of `Style` if the current user is root.
	// If this style is empty, will fall back to Style.
	RootStyle string `yaml:"rootStyle"`
	// ViCmdPrompt is what to display as the prompt if the shell is in vicmd mode.
	// Defaults to ": ".
	VicmdPrompt string `yaml:"vicmdPrompt"`
	// VicmdStyle will be used when the shell is in vicmd mode.
	VicmdStyle string `yaml:"vicmdStyle"`
	// ErrorStyle will be used when the previous command failed.
	ErrorStyle string `yaml:"errorStyle"`
}

type promptModuleData struct {
	// PromptString is the chosen prompt string, before styling.
	PromptString string
	// PromptStyle is the chosen prompt style.
	PromptStyle string
	// ViCmdMode is true if the shell is in vicmd mode.
	ViCmdMode bool
}

// Execute the prompt module.
func (mod PromptModule) Execute(context *Context) ModuleResult {
	var text string
	var style string

	viCmdMode := context.Globals.Keymap == "vicmd"

	if viCmdMode {
		text = mod.VicmdPrompt
		style = defaultString(mod.VicmdStyle, mod.Style)
	} else if !context.Globals.IsRoot {
		text = mod.Prompt
		style = mod.Style
	} else {
		text = mod.RootPrompt
		style = defaultString(mod.RootStyle, mod.Style)
	}

	if context.Globals.Status != 0 {
		style = defaultString(mod.ErrorStyle, mod.Style)
	}

	data := promptModuleData{
		PromptString: text,
		PromptStyle:  style,
		ViCmdMode:    viCmdMode,
	}

	return executeModule(context, mod.CommonConfig, data, style, text)
}

func init() {
	registerModule(
		"prompt",
		registeredModule{
			jsonSchema: schemas.PromptModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := PromptModule{
					Type:        "prompt",
					Prompt:      "$ ",
					RootPrompt:  "# ",
					VicmdPrompt: ": ",
				}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
