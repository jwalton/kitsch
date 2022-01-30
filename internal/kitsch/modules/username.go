package modules

import (
	"os/user"

	"github.com/jwalton/kitsch/internal/kitsch/log"
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas UsernameModule

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
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=username"`
	// ShowAlways will cause the username to always be shown.  If false (the default),
	// then the username will only be shown if the user is root, or the current
	// session is an SSH session.
	ShowAlways bool `yaml:"showAlways"`
	// RootStyle will be used in place of `Style` if the current user is root.
	// If this style is empty, will fall back to `Style`.
	RootStyle string `yaml:"rootStyle"`
}

type usernameModuleData struct {
	// username is the current user's username.
	username string
	// IsSSH is true if the user is in an SSH session.
	IsSSH bool
	// Show is true if the username module should be displayed.
	Show bool
}

// Username is the current user's username.
func (data usernameModuleData) Username() string {
	if data.username != "" {
		return data.username
	}

	// Fetch the user from the OS.  This can be a little slow, eating up around
	// 6ms on MaxOS and Linux style systems, which is why we prefer to get
	// the username from the env.  The good news is that `os/user` caches this
	// value for us, so repeated calls shouldn't be slow.
	user, err := user.Current()
	if err != nil {
		log.Info("Unable to get current user: " + err.Error())
		return ""
	}
	return user.Username
}

// Execute the username module.
func (mod UsernameModule) Execute(context *Context) ModuleResult {
	isRoot := context.Globals.IsRoot
	isSSH := context.Environment.HasSomeEnv("SSH_CLIENT", "SSH_CONNECTION", "SSH_TTY")
	show := isSSH || isRoot || mod.ShowAlways

	data := usernameModuleData{
		username: context.Environment.Getenv("USER"),
		IsSSH:    isSSH,
		Show:     show,
	}

	defaultText := ""
	style := ""

	if show {
		defaultText = data.Username()
		if isRoot && mod.RootStyle != "" {
			style = mod.RootStyle
		}
	}

	return ModuleResult{
		DefaultText:   defaultText,
		StyleOverride: style,
		Data:          data,
	}
}

func init() {
	registerModule(
		"username",
		registeredModule{
			jsonSchema: schemas.UsernameModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := UsernameModule{Type: "username"}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
