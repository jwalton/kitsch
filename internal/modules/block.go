package modules

import (
	"text/template"

	"github.com/jwalton/kitsch-prompt/internal/env"
)

// BlockConfig is configuration for a block module.
type BlockConfig struct {
	CommonConfig
	// Modules is a list of child modules to be rendered under this block
	Modules []Module
	// Join is a template to use to join together modules.  This will be executed
	// with the following parameters:
	//
	// • prevStyle - The style of last character of the previous module.
	//
	// • nextStyle - The style of the first character of the next module.
	//
	// • index - The index of the next module in the Modules array.
	//
	Join string
}

type block struct {
	config BlockConfig
}

// NewBlockModule creates a "block" module, which renders a collection of other
// modules.
//
// Any module that outputs no text is considered "inactive" and will not be
// part of the result.
//
// The block module returns the following template variables:
//
// • children - The results of executing each child module.  Only modules that
//   actually generated output will be included.
//
func NewBlockModule(config BlockConfig) Module {
	return block{config}
}

func (mod block) Execute(env env.Env) ModuleResult {
	config := mod.config

	children := make([]ModuleResult, 0, len(mod.config.Modules))
	for _, module := range mod.config.Modules {
		result := module.Execute(env)
		if len(result.Text) != 0 {
			children = append(children, result)
		}
	}

	defaultText := mod.joinChildren(children)

	data := map[string]interface{}{
		"children":    children,
		"defaultText": defaultText,
	}

	result := executeModule(config.CommonConfig, data, config.Style, defaultText)

	if len(children) > 0 {
		result.StartStyle = config.Style.Mix(children[0].StartStyle)
		result.EndStyle = config.Style.Mix(children[len(children)-1].EndStyle)
	}

	return result
}

func (mod block) joinChildren(children []ModuleResult) string {
	result := ""
	var join *template.Template = nil

	if mod.config.Join != "" {
		var err error
		join, err = template.New("join").Parse(mod.config.Join)
		if err != nil {
			join = nil
		}
	}

	for index, child := range children {
		if join != nil && index != 0 {
			prev := children[index-1]
			result += templateToString(join, map[string]interface{}{
				"prevStyle": prev.EndStyle,
				"nextStyle": child.StartStyle,
				"index":     index,
			})
		}

		result += child.Text
	}

	return result
}
