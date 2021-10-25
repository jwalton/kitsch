package modules

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// HostnameModule shows the name of the current hostname.  This is,
// by default, hidden unless the session is an SSH session.
//
// The hostname module provides the following template variables:
//
// • Hostname - The current hostname.
//
// • IsSSH - True if this is an SSH session, false otherwise.
//
// • Show - True if we should show the hostname module, false otherwise.
//
type HostnameModule struct {
	CommonConfig `yaml:",inline"`
	// ShowAlways will cause the hostname to always be shown.  If false (the default),
	// then the hostname will only be shown if the current session is an SSH session.
	ShowAlways bool `yaml:"showAlways"`
}

// Execute the module.
func (mod HostnameModule) Execute(context *Context) ModuleResult {
	// TODO: Move isSSH to somewhere common.
	isSSH := context.Environment.HasSomeEnv("SSH_CLIENT", "SSH_CONNECTION", "SSH_TTY")
	show := isSSH || mod.ShowAlways

	hostname := context.Globals.Hostname

	// If the hostname is a FQDM, just grab the first part of the hostname.
	if strings.Contains(hostname, ".") {
		hostname = strings.Split(hostname, ".")[0]
	}

	data := map[string]interface{}{
		"Hostname": hostname,
		"IsSSH":    isSSH,
		"Show":     show,
	}

	defaultText := ""

	if show {
		defaultText = context.Globals.Username
	}

	return executeModule(context, mod.CommonConfig, data, mod.Style, defaultText)
}

func init() {
	registerFactory("hostname", func(node *yaml.Node) (Module, error) {
		var module HostnameModule
		err := node.Decode(&module)
		return &module, err
	})
}
