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
	// Prompt is what to display as the prompt.  Defaults to "$".
	Prompt string `yaml:"prompt"`
	// RootPrompt is what to display as the prompt if the current user is root.  Defaults to "#".
	RootPrompt string `yaml:"rootPrompt"`
	// RootStyle will be used in place of `Style` if the current user is root.
	// If this style is empty, will fall back to Style.
	RootStyle string `yaml:"rootStyle"`
	// ErrorStyle will be used when the previous command failed.
	ErrorStyle string `yaml:"errorStyle"`
}

// Execute the prompt module.
func (mod PromptModule) Execute(context *Context) ModuleResult {
	isRoot := context.Environment.IsRoot()

	data := map[string]interface{}{
		"IsRoot": isRoot,
	}

	// TODO: Use globals.Keymap here.
	var text string
	if !isRoot {
		text = defaultString(mod.Prompt, "$")
	} else {
		text = defaultString(mod.RootPrompt, "#")
	}

	var style string
	if context.Globals.Status != 0 {
		style = defaultString(mod.ErrorStyle, mod.Style)
	} else if !isRoot {
		style = mod.Style
	} else {
		style = defaultString(mod.RootStyle, mod.Style)
	}

	return executeModule(context, mod.CommonConfig, data, style, text+" ")
}

func init() {
	registerFactory("prompt", func(node *yaml.Node) (Module, error) {
		var module PromptModule
		err := node.Decode(&module)
		return &module, err
	})
}
