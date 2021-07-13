package modules

import (
	"github.com/jwalton/kitsch-prompt/internal/env"
	styleLib "github.com/jwalton/kitsch-prompt/internal/style"
)

// PromptConfig is configuration for a prompt module.
type PromptConfig struct {
	CommonConfig
	// Prompt is what to display as the prompt.  Defaults to "$".
	Prompt string
	// RootPrompt is what to display as the prompt if the current user is root.  Defaults to "#".
	RootPrompt string
	// RootStyle will be used in place of `Style` if the current user is root.
	// If this style is empty, will fall back to Style.
	RootStyle styleLib.Style
	// ErrorStyle will be used when the previous command failed.
	ErrorStyle styleLib.Style
}

type prompt struct {
	config PromptConfig
}

// NewPromptModule creates a prompt module.
//
// The prompt module displays a "$", or a "#" if the current user is root.
//
// The prompt module returns the following template variables:
//
// • isRoot - True if the user is root, false otherwise.
//
// • status - An `int` which represents the return status of the last command.
//
func NewPromptModule(config PromptConfig) Module {
	return prompt{config}
}

func (mod prompt) Execute(env env.Env) ModuleResult {
	config := mod.config

	isRoot := env.IsRoot()
	status := env.Status()

	data := map[string]interface{}{
		"isRoot": isRoot,
		"status": status,
	}

	// TODO: Use env.Keymap() here.
	var text string
	if !isRoot {
		text = defaultString(config.Prompt, "$")
	} else {
		text = defaultString(config.RootPrompt, "#")
	}

	var style styleLib.Style
	if status != 0 {
		style = config.ErrorStyle
	} else if !isRoot {
		style = config.Style
		if style.IsEmpty() {
			style = config.Style
		}
	} else {
		style = config.RootStyle
		if style.IsEmpty() {
			style = config.Style
		}
	}

	return executeModule(config.CommonConfig, data, style, text)
}
