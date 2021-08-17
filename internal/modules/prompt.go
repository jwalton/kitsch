package modules

import (
	"github.com/jwalton/kitsch-prompt/internal/env"
	styleLib "github.com/jwalton/kitsch-prompt/internal/style"
	"gopkg.in/yaml.v3"
)

// PromptModule shows a prompt to the user.
//
// The prompt module displays a "$", or a "#" if the current user is root.
//
// The prompt module provides the following template variables:
//
// • isRoot - True if the user is root, false otherwise.
//
// • status - An `int` which represents the return status of the last command.
//
type PromptModule struct {
	CommonConfig `yaml:",inline"`
	// Prompt is what to display as the prompt.  Defaults to "$".
	Prompt string
	// RootPrompt is what to display as the prompt if the current user is root.  Defaults to "#".
	RootPrompt string
	// RootStyle will be used in place of `Style` if the current user is root.
	// If this style is empty, will fall back to Style.
	RootStyle styleLib.Style `yaml:"rootStyle"`
	// ErrorStyle will be used when the previous command failed.
	ErrorStyle styleLib.Style `yaml:"errorStyle"`
}

// Execute the prompt module.
func (mod PromptModule) Execute(env env.Env) ModuleResult {
	isRoot := env.IsRoot()
	status := env.Status()

	data := map[string]interface{}{
		"isRoot": isRoot,
		"status": status,
	}

	// TODO: Use env.Keymap() here.
	var text string
	if !isRoot {
		text = defaultString(mod.Prompt, "$")
	} else {
		text = defaultString(mod.RootPrompt, "#")
	}

	var style styleLib.Style
	if status != 0 {
		style = mod.ErrorStyle
		if style.IsEmpty() {
			style = mod.Style
		}
	} else if !isRoot {
		style = mod.Style
	} else {
		style = mod.RootStyle
		if style.IsEmpty() {
			style = mod.Style
		}
	}

	return executeModule(mod.CommonConfig, data, style, text)
}

func init() {
	registerFactory("prompt", func(node *yaml.Node) (Module, error) {
		var module PromptModule
		err := node.Decode(&module)
		return &module, err
	})
}
