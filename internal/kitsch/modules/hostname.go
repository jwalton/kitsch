package modules

import (
	"strings"

	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas HostnameModule

// HostnameModule shows the name of the current hostname.  This is,
// by default, hidden unless the session is an SSH session.
//
type HostnameModule struct {
	CommonConfig `yaml:",inline"`
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",enum=hostname"`
	// ShowAlways will cause the hostname to always be shown.  If false (the default),
	// then the hostname will only be shown if the current session is an SSH session.
	ShowAlways bool `yaml:"showAlways"`
}

type hostnameResult struct {
	// Hostname is the current hostname.
	Hostname string `yaml:"hostname"`
	// IsSSH is true if this is an SSH session, false otherwise.
	IsSSH bool `yaml:"isSSH"`
	// Show is true if we should show the hostname, false otherwise.
	Show bool `yaml:"show"`
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

	data := hostnameResult{
		Hostname: hostname,
		IsSSH:    isSSH,
		Show:     show,
	}

	defaultText := ""
	if show {
		defaultText = hostname
	}

	return executeModule(context, mod.CommonConfig, data, mod.Style, defaultText)
}

func init() {
	registerModule(
		"hostname",
		registeredModule{
			jsonSchema: schemas.HostnameModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := HostnameModule{Type: "hostname"}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
