package modules

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// CmdDurationModule shows the amount of time the previous command took to execute.
//
// The module provides the following template variables:
//
// • Duration - A string describing the duration of the previous command.
//   Defaults to 2000ms.
//
type CmdDurationModule struct {
	CommonConfig `yaml:",inline"`
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

func init() {
	registerFactory("command_duration", func(node *yaml.Node) (Module, error) {
		module := CmdDurationModule{MinTime: 2000}
		err := node.Decode(&module)
		return &module, err
	})
}
