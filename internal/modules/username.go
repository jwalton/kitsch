package modules

import (
	"github.com/jwalton/kitsch-prompt/internal/env"
	styleLib "github.com/jwalton/kitsch-prompt/internal/style"
	"gopkg.in/yaml.v3"
)

// The UsernameModule shows the name of the currently logged in user.  This is,
// by default, hidden unless the user is root or the session is an SSH session.
// The CommonConfig.Style is applied by default, unless the user is Root in which
// case it is overridden by `UsernameConfig.RootStyle`.
//
// The username module provides the following template variables:
//
// • username - The current user's username.
//
// • isRoot - True if the user is root, false otherwise.
//
// • isSSH - True if this is an SSH session, false otherwise.
//
// • show - True if we should show the username module, false otherwise.
//
type UsernameModule struct {
	CommonConfig `yaml:",inline"`
	// ShowAlways will cause the username to always be shown.  If false (the default),
	// then the username will only be shown if the user is root, or the current
	// session is an SSH session.
	ShowAlways bool
	// RootStyle will be used in place of `Style` if the current user is root.
	// If this style is empty, will fall back to `Style`.
	RootStyle styleLib.Style
}

// Execute the username module.
func (mod UsernameModule) Execute(env env.Env) ModuleResult {
	username := env.GetUsername()
	isRoot := env.IsRoot()
	isSSH := env.HasSomeEnv("SSH_CLIENT", "SSH_CONNECTION", "SSH_TTY")
	show := isSSH || isRoot || mod.ShowAlways

	data := map[string]interface{}{
		"Username": username,
		"IsRoot":   isRoot,
		"IsSSH":    isSSH,
		"Show":     show,
	}

	defaultText := ""
	var style styleLib.Style

	if show {
		defaultText = username
		if isRoot && !mod.RootStyle.IsEmpty() {
			style = mod.RootStyle
		} else {
			style = mod.Style
		}
	}

	return executeModule(mod.CommonConfig, data, style, defaultText)
}

func init() {
	registerFactory("username", func(node *yaml.Node) (Module, error) {
		var module UsernameModule
		err := node.Decode(&module)
		return &module, err
	})
}
