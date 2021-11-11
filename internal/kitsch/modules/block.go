package modules

import (
	"strings"
	"text/template"
	"time"

	"github.com/jwalton/kitsch-prompt/internal/kitsch/modtemplate"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/styling"
	"gopkg.in/yaml.v3"
)

// BlockModule renders a collection of other modules.
//
// Any module that outputs no text is considered "inactive" and will not be
// part of the result.
//
// Provides the following template variables:
//
// • ModuleArray - The results of executing each child module.  Only modules that
//   actually generated output will be included.
//
// • Modules - A map of results from executing each child module, indexed by
//   module ID.  Only modules that actually generated output will be included.
//
type BlockModule struct {
	CommonConfig `yaml:",inline"`
	// Modules is a list of child modules to be rendered under this block
	Modules []ModuleSpec `yaml:"modules"`
	// Join is a template to use to join together modules.  Defaults to " ".
	// This will be executed with the following parameters:
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
	resultsArray := make([]ModuleResult, 0, len(mod.Modules))
	childDurations := make([]ModuleDuration, 0, len(mod.Modules))
	resultsByID := make(map[string]ModuleResult, len(mod.Modules))

	for index := range mod.Modules {
		item := &mod.Modules[index]
		start := time.Now()

		result := item.Module.Execute(context)

		childDurations = append(childDurations, ModuleDuration{
			Module:   item,
			Duration: time.Since(start),
			Children: result.ChildDurations,
		})

		if len(result.Text) != 0 {
			resultsArray = append(resultsArray, result)
			if item.ID != "" {
				resultsByID[item.ID] = result
			}
		}
	}

	defaultText := mod.joinChildren(context, resultsArray)

	data := map[string]interface{}{
		"Modules":     resultsByID,
		"ModuleArray": resultsArray,
	}

	result := executeModule(context, mod.CommonConfig, data, mod.Style, defaultText)

	if len(resultsArray) > 0 {
		lastChild := len(resultsArray) - 1
		result.StartStyle = styling.CharacterColors{
			FG: defaultString(result.StartStyle.FG, resultsArray[0].StartStyle.FG),
			BG: defaultString(result.StartStyle.BG, resultsArray[0].StartStyle.BG),
		}
		result.EndStyle = styling.CharacterColors{
			FG: defaultString(result.EndStyle.FG, resultsArray[lastChild].EndStyle.FG),
			BG: defaultString(result.EndStyle.BG, resultsArray[lastChild].EndStyle.BG),
		}
	}

	result.ChildDurations = childDurations

	return result
}

// blockJoinData is the data passed to the join template.
type blockJoinData struct {
	Global     *Globals
	PrevColors styling.CharacterColors
	NextColors styling.CharacterColors
	Index      int
}

func (mod BlockModule) joinChildren(context *Context, children []ModuleResult) string {
	out := strings.Builder{}

	var join *template.Template = nil

	if !strings.Contains(mod.Join, "{{") {
		// Not a template, just a string.
		for index, child := range children {
			if index != 0 {
				out.WriteString(mod.Join)
			}
			out.WriteString(child.Text)
		}

	} else {
		// Compile the join template
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
				out.WriteString(joiner)
			}

			out.WriteString(child.Text)
		}
	}

	return out.String()
}

func init() {
	registerFactory("block", func(node *yaml.Node) (Module, error) {
		module := BlockModule{Join: " "}
		err := node.Decode(&module)
		return &module, err
	})
}
