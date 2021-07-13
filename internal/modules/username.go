package modules

import (
	"github.com/jwalton/kitsch-prompt/internal/env"
	styleLib "github.com/jwalton/kitsch-prompt/internal/style"
)

// UsernameConfig is configuration for a username module.
type UsernameConfig struct {
	CommonConfig
	// ShowAlways will cause the username to always be shown.  If false (the default),
	// then the username will only be shown if the user is root, or the current
	// session is an SSH session.
	ShowAlways bool
	// RootStyle will be used in place of `Style` if the current user is root.
	// If this style is empty, will fall back to `Style`.
	RootStyle styleLib.Style
}

type username struct {
	config UsernameConfig
}

// NewUsernameModule creates a username module.
//
// The UsernameModule shows the name of the currently logged in user.  This is,
// by default, hidden unless the user is root or the session is an SSH session.
// The CommonConfig.Style is applied by default, unless the user is Root in which
// case it is overridden by `UsernameConfig.RootStyle`.
//
// The username module returns the following template variables:
//
// • username - The current user's username.
//
// • isRoot - True if the user is root, false otherwise.
//
// • isSSH - True if this is an SSH session, false otherwise.
//
// • show - True if we should show the username module, false otherwise.
//
func NewUsernameModule(config UsernameConfig) Module {
	return username{config}
}

func (mod username) Execute(env env.Env) ModuleResult {
	config := mod.config

	username := env.GetUsername()
	isRoot := env.IsRoot()
	isSSH := env.HasSomeEnv("SSH_CLIENT", "SSH_CONNECTION", "SSH_TTY")
	show := isSSH || isRoot || config.ShowAlways

	data := map[string]interface{}{
		"config":   config,
		"username": username,
		"isRoot":   isRoot,
		"isSSH":    isSSH,
		"show":     show,
	}

	defaultText := ""
	var style styleLib.Style

	if show {
		defaultText = username
		if isRoot && !config.RootStyle.IsEmpty() {
			style = config.RootStyle
		} else {
			style = config.Style
		}
	}

	return executeModule(config.CommonConfig, data, style, defaultText)
}
