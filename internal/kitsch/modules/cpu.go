package modules

import (
	"fmt"
	"time"

	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"github.com/shirou/gopsutil/v3/cpu"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas CPUModule

// CPUModule shows the current CPU usage.
type CPUModule struct {
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=cpu"`
	// MinPercent is the minimum percentage to show the CPU.  Defaults to 5%.
	MinPercent float64 `yaml:"minPercent"`
}

type cpuModuleData struct {
	// Percent is the overall CPU usage as a percentage.
	Percent float64
}

// Execute the time module.
func (mod CPUModule) Execute(context *Context) ModuleResult {
	result := cpuModuleData{}

	// TODO: This needs a positive delay to return meaningful results.
	// Need to move this into a background task.
	percents, err := cpu.Percent(100*time.Millisecond, false)
	if err == nil && len(percents) > 0 {
		result.Percent = percents[0]
	}

	defaultText := ""
	if result.Percent > mod.MinPercent {
		defaultText = fmt.Sprintf("%.0f%%", result.Percent)
	}

	return ModuleResult{
		DefaultText: defaultText,
		Data:        result,
	}
}

func init() {
	registerModule(
		"cpu",
		registeredModule{
			jsonSchema: schemas.CPUModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := CPUModule{Type: "cpu", MinPercent: 5}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
