package modules

import (
	"text/template"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/modtemplate"
	"gopkg.in/yaml.v3"
)

// BlockModule renders a collection of other modules.
//
// Any module that outputs no text is considered "inactive" and will not be
// part of the result.
//
// Provides the following template variables:
//
// • children - The results of executing each child module.  Only modules that
//   actually generated output will be included.
//
type BlockModule struct {
	CommonConfig `yaml:",inline"`
	// Modules is a list of child modules to be rendered under this block
	Modules ModuleList
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

// Execute the block module.
func (mod BlockModule) Execute(env env.Env) ModuleResult {
	children := make([]ModuleResult, 0, len(mod.Modules.Modules))
	for _, module := range mod.Modules.Modules {
		result := module.Execute(env)
		if len(result.Text) != 0 {
			children = append(children, result)
		}
	}

	defaultText := mod.joinChildren(children)

	data := map[string]interface{}{
		"children": children,
	}

	result := executeModule(mod.CommonConfig, data, mod.Style, defaultText)

	if len(children) > 0 {
		result.StartStyle = mod.Style.Mix(children[0].StartStyle)
		result.EndStyle = mod.Style.Mix(children[len(children)-1].EndStyle)
	}

	return result
}

func (mod BlockModule) joinChildren(children []ModuleResult) string {
	result := ""
	var join *template.Template = nil

	if mod.Join != "" {
		var err error
		join, err = modtemplate.CompileTemplate("join", mod.Join)
		if err != nil {
			join = nil
		}
	}

	for index, child := range children {
		if join != nil && index != 0 {
			prev := children[index-1]
			joiner, err := modtemplate.TemplateToString(join, map[string]interface{}{
				"prevStyle": prev.EndStyle,
				"nextStyle": child.StartStyle,
				"index":     index,
			})
			if err != nil {
				// TODO: Add warning
				joiner = " "
			}
			result += joiner
		}

		result += child.Text
	}

	return result
}

func init() {
	registerFactory("block", func(node *yaml.Node) (Module, error) {
		var module BlockModule
		err := node.Decode(&module)
		return &module, err
	})
}
