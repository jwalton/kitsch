package modules

import (
	"gopkg.in/yaml.v3"
)

// UsernameModule shows the name of the currently logged in user.  This is,
// by default, hidden unless the user is root or the session is an SSH session.
// The CommonConfig.Style is applied by default, unless the user is Root in which
// case it is overridden by `UsernameConfig.RootStyle`.
//
// The username module provides the following template variables:
//
// • Username - The current user's username.
//
// • IsRoot - True if the user is root, false otherwise.
//
// • IsSSH - True if this is an SSH session, false otherwise.
//
// • Show - True if we should show the username module, false otherwise.
//
type UsernameModule struct {
	CommonConfig `yaml:",inline"`
	// ShowAlways will cause the username to always be shown.  If false (the default),
	// then the username will only be shown if the user is root, or the current
	// session is an SSH session.
	ShowAlways bool `yaml:"showAlways"`
	// RootStyle will be used in place of `Style` if the current user is root.
	// If this style is empty, will fall back to `Style`.
	RootStyle string `yaml:"rootStyle"`
}

type usernameModuleData struct {
	// Username is the current user's username.
	Username string
	// FullName is the current user's full name, if available.
	FullName string
	// IsRoot is true if the current user is root.
	IsRoot bool
	// IsSSH is true if the user is in an SSH session.
	IsSSH bool
	// Show is true if the username module should be displayed.
	Show bool
}

// Execute the username module.
func (mod UsernameModule) Execute(context *Context) ModuleResult {
	isRoot := context.Environment.IsRoot()
	isSSH := context.Environment.HasSomeEnv("SSH_CLIENT", "SSH_CONNECTION", "SSH_TTY")
	show := isSSH || isRoot || mod.ShowAlways

	data := usernameModuleData{
		Username: context.Globals.Username,
		FullName: context.Globals.UserFullName,
		IsRoot:   isRoot,
		IsSSH:    isSSH,
		Show:     show,
	}

	defaultText := ""
	style := mod.Style

	if show {
		defaultText = context.Globals.Username
		if isRoot && mod.RootStyle != "" {
			style = mod.RootStyle
		}
	}

	return executeModule(context, mod.CommonConfig, data, style, defaultText)
}

func init() {
	registerFactory("username", func(node *yaml.Node) (Module, error) {
		var module UsernameModule
		err := node.Decode(&module)
		return &module, err
	})
}
