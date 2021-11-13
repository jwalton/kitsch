package modules

import (
	"gopkg.in/yaml.v3"
)

// PromptModule shows a prompt to the user.
//
// The prompt module displays a "$", or a "#" if the current user is root.
//
// The prompt module provides the following template variables:
//
// • IsRoot - True if the user is root, false otherwise.
//
// • Status - An `int` which represents the return status of the last command.
//
type PromptModule struct {
	CommonConfig `yaml:",inline"`
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

// Execute the prompt module.
func (mod PromptModule) Execute(context *Context) ModuleResult {
	isRoot := context.Environment.IsRoot()

	data := map[string]interface{}{
		"IsRoot": isRoot,
	}

	var text string
	var style string
	if context.Globals.Keymap == "vicmd" {
		text = mod.VicmdPrompt
		style = defaultString(mod.VicmdStyle, mod.Style)
	} else if !isRoot {
		text = mod.Prompt
		style = mod.Style
	} else {
		text = mod.RootPrompt
		style = defaultString(mod.RootStyle, mod.Style)
	}

	if context.Globals.Status != 0 {
		style = defaultString(mod.ErrorStyle, mod.Style)
	}

	return executeModule(context, mod.CommonConfig, data, style, text)
}

func init() {
	registerFactory("prompt", func(node *yaml.Node) (Module, error) {
		module := PromptModule{
			Prompt:      "$ ",
			RootPrompt:  "# ",
			VicmdPrompt: ": ",
		}
		err := node.Decode(&module)
		return &module, err
	})
}
