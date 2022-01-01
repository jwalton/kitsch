package modules

import (
	"fmt"
	"time"

	"github.com/jwalton/kitsch-prompt/internal/kitsch/log"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas CmdDurationModule

// CmdDurationModule shows the amount of time the previous command took to execute.
//
// The module provides the following template variables:
//
// • Duration - A string describing the duration of the previous command.
//   Defaults to 2000ms.
//
type CmdDurationModule struct {
	CommonConfig `yaml:",inline"`
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",enum=command_duration"`
	// MinTime is the minimum duration to show, in milliseconds.
	MinTime int64 `yaml:"minTime"`
	// ShowMilliseconds - If true, show milliseconds.
	ShowMilliseconds bool `yaml:"showMilliseconds"`
}

type cmdDurationModuleResult struct {
	// Duration is the duration the command took, in milliseconds.
	Duration int64
	// PrettyDuration is the duration the command took, in a human-readable format.
	PrettyDuration string
}

// Execute the module.
func (mod CmdDurationModule) Execute(context *Context) ModuleResult {
	var durationStr string
	if context.Globals.PreviousCommandDuration < mod.MinTime {
		durationStr = ""
	} else {
		durationStr = mod.formatDuration(context.Globals.PreviousCommandDuration)
	}

	data := cmdDurationModuleResult{
		Duration:       context.Globals.PreviousCommandDuration,
		PrettyDuration: durationStr,
	}

	return executeModule(context, mod.CommonConfig, data, mod.Style, durationStr)
}

func (mod CmdDurationModule) formatDuration(timeInMs int64) string {
	d := time.Duration(timeInMs) * time.Millisecond
	result := d.Round(time.Second).String()

	if mod.ShowMilliseconds {
		result = fmt.Sprintf("%s%dms", result, timeInMs%1000)
	}

	return result
}

// Validate validates this module.
func (mod CmdDurationModule) Validate(context *Context, prefix string) {
	mod.CommonConfig.validate(context, prefix)
	if mod.MinTime < 0 {
		log.Warn(fmt.Sprintf("%s: Invalid minTime: %d", prefix, mod.MinTime))
	}

	testTemplate(context, prefix, mod.Template, map[string]interface{}{
		"Zero duration": cmdDurationModuleResult{Duration: 0, PrettyDuration: ""},
		"With duration": cmdDurationModuleResult{Duration: 20, PrettyDuration: "20s"},
	})
}

func init() {
	registerModule(
		"command_duration",
		registeredModule{
			jsonSchema: schemas.CmdDurationModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := CmdDurationModule{Type: "command_duration", MinTime: 2000}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
