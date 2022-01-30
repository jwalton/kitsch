package modules

import (
	"time"

	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas TimeModule

const defaultTimeFormat = "15:04:05"

// TimeModule shows the current time.
type TimeModule struct {
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=time"`
	// Layout is the format to show the time in.  Layout defines the format by
	// showing how the reference time, defined to be
	//
	//     Mon Jan 2 15:04:05 -0700 MST 2006
	//
	// (See https://golang.org/pkg/time/#Time.Format for more details.)
	//
	// Defaults to "15:04:05".
	//
	Layout string `yaml:"layout"`
}

type timeModuleData struct {
	// Time is the current time, as a `time.Time` object.
	Time time.Time
	// Unix is the number of seconds since the Unix epoch.
	Unix int64
	// TimeStr is the current time as a formatted string.
	TimeStr string
}

// Execute the time module.
func (mod TimeModule) Execute(context *Context) ModuleResult {
	now := time.Now()

	layout := mod.Layout
	if layout == "" {
		layout = defaultTimeFormat
	}

	formattedTime := now.Format(layout)

	return ModuleResult{
		DefaultText: formattedTime,
		Data: timeModuleData{
			Time:    now,
			Unix:    now.Unix(),
			TimeStr: formattedTime,
		},
	}
}

func init() {
	registerModule(
		"time",
		registeredModule{
			jsonSchema: schemas.TimeModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := TimeModule{Type: "time"}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
