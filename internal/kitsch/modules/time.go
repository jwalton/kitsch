package modules

import (
	"time"

	"gopkg.in/yaml.v3"
)

const defaultTimeFormat = "15:04:05"

// TimeModule shows the current time.
//
// Provides the following template variables:
//
// • time - The current time, as a `time.Time` object.
//
// • timeStr - The current time, as a formatted string.
//
type TimeModule struct {
	CommonConfig `yaml:",inline"`
	// Layout is the format to show the time in.  Layout defines the format by
	// showing how the reference time, defined to be
	//
	//     Mon Jan 2 15:04:05 -0700 MST 2006
	//
	// (See https://golang.org/pkg/time/#Time.Format for more details.)
	//
	// Defaults to "15:04:05".
	//
	Layout string
}

// Execute the time module.
func (mod TimeModule) Execute(context *Context) ModuleResult {
	now := time.Now()

	layout := mod.Layout
	if layout == "" {
		layout = defaultTimeFormat
	}

	formattedTime := now.Format(layout)

	data := map[string]interface{}{
		"Time":    now,
		"TimeStr": formattedTime,
	}

	return executeModule(context, mod.CommonConfig, data, mod.Style, formattedTime)
}

func init() {
	registerFactory("time", func(node *yaml.Node) (Module, error) {
		var module TimeModule
		err := node.Decode(&module)
		return &module, err
	})
}
