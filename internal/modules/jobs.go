package modules

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// JobsModule shows the current count of running background jobs.  If
// the number of running jobs is greater than or equal to "SymbolThreshold",
// then the "Symbol" will be shown.  If the number of running jobs is greater
// than or equal to "CountThreshold", then the count of running jobs will be
// shown.
//
type JobsModule struct {
	CommonConfig `yaml:",inline"`
	// Symbol is the symbol to show when there are background jobs.  Defaults to "+".
	Symbol string `yaml:"symbol"`
	// SymbolThreshold is the threshold for showing the symbol.  Defaults to 1.
	SymbolThreshold int `yaml:"symbolThreshold"`
	// CountThreshold is the threshold for showing the count of background jobs.  Defaults to 2.
	CountThreshold int `yaml:"countThreshold"`
}

type jobsModuleData struct {
	// Jobs is the count of running jobs.
	Jobs int
	// ShowSymbol is true if the symbol should be shown.
	ShowSymbol bool
	// ShowCount is true if the count should be shown.
	ShowCount bool
}

// Execute the module.
func (mod JobsModule) Execute(context *Context) ModuleResult {
	jobs := context.Environment.Jobs()
	showSymbol := jobs >= mod.SymbolThreshold
	showCount := jobs >= mod.CountThreshold

	data := jobsModuleData{
		Jobs:       jobs,
		ShowSymbol: showSymbol,
		ShowCount:  showCount,
	}

	defaultText := ""

	if showSymbol {
		defaultText += mod.Symbol
	}
	if showCount {
		defaultText += fmt.Sprintf("%d", jobs)
	}

	return executeModule(context, mod.CommonConfig, data, mod.Style, defaultText)
}

func init() {
	registerFactory("jobs", func(node *yaml.Node) (Module, error) {
		module := JobsModule{
			Symbol:          "+",
			SymbolThreshold: 1,
			CountThreshold:  2,
		}
		err := node.Decode(&module)
		return &module, err
	})
}
