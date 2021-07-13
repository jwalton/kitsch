package modules

import (
	"time"

	"github.com/jwalton/kitsch-prompt/internal/env"
)

const defaultTimeFormat = "15:04:05"

// TimeConfig is configuration for a time module.
type TimeConfig struct {
	CommonConfig
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

type timeModule struct {
	config TimeConfig
}

// NewTimeModule creates a time module.
//
// Returns the following template variables:
//
// • time - The current time, as a `time.Time` object.
//
// • timeStr - The current time, as a formatted string.
//
func NewTimeModule(config TimeConfig) Module {
	return timeModule{config}
}

func (mod timeModule) Execute(env env.Env) ModuleResult {
	config := mod.config
	now := time.Now()

	layout := config.Layout
	if layout == "" {
		layout = defaultTimeFormat
	}

	formattedTime := now.Format(layout)

	data := map[string]interface{}{
		"time":    now,
		"timeStr": formattedTime,
	}

	return executeModule(config.CommonConfig, data, config.Style, formattedTime)
}
