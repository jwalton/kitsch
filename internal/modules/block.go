package modules

import (
	"text/template"

	"github.com/jwalton/kitsch-prompt/internal/modtemplate"
	"github.com/jwalton/kitsch-prompt/internal/style"
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
	// • PrevColors - The FG and BG color of last character of the previous module.
	//
	// • NextColors - The FG and BG of the first character of the next module.
	//
	// • Index - The index of the next module in the Modules array.
	//
	Join string
}

// Execute the block module.
func (mod BlockModule) Execute(context *Context) ModuleResult {
	children := make([]ModuleResult, 0, len(mod.Modules.Modules))
	for _, module := range mod.Modules.Modules {
		result := module.Execute(context)
		if len(result.Text) != 0 {
			children = append(children, result)
		}
	}

	defaultText := mod.joinChildren(context, children)

	data := map[string]interface{}{
		"Children": children,
	}

	result := executeModule(context, mod.CommonConfig, data, mod.Style, defaultText)

	if len(children) > 0 {
		lastChild := len(children) - 1
		result.StartStyle = style.CharacterColors{
			FG: defaultString(result.StartStyle.FG, children[0].StartStyle.FG),
			BG: defaultString(result.StartStyle.BG, children[0].StartStyle.BG),
		}
		result.EndStyle = style.CharacterColors{
			FG: defaultString(result.EndStyle.FG, children[lastChild].EndStyle.FG),
			BG: defaultString(result.EndStyle.BG, children[lastChild].EndStyle.BG),
		}
	}

	return result
}

// blockJoinData is the data passed to the join template.
type blockJoinData struct {
	Global     *Globals
	PrevColors style.CharacterColors
	NextColors style.CharacterColors
	Index      int
}

func (mod BlockModule) joinChildren(context *Context, children []ModuleResult) string {
	result := ""
	var join *template.Template = nil

	if mod.Join != "" {
		var err error
		join, err = modtemplate.CompileTemplate(&context.Styles, "join", mod.Join)
		if err != nil {
			join = nil
		}
	}

	for index, child := range children {
		if join != nil && index != 0 {
			prev := children[index-1]
			joiner, err := modtemplate.TemplateToString(join, blockJoinData{
				Global:     &context.Globals,
				PrevColors: prev.EndStyle,
				NextColors: child.StartStyle,
				Index:      index,
			})
			if err != nil {
				context.Environment.Warn(err.Error())
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
