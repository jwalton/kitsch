package modules

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jwalton/kitsch-prompt/internal/kitsch/log"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/modtemplate"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/styling"
	"github.com/jwalton/kitsch-prompt/internal/perf"
	"gopkg.in/yaml.v3"
)

var blockSchema = `{
	"type":"object",
	"properties": {
		"type": {
			"enum": [ "block" ],
			"type": "string"
		},
	 	"style": {"type": "string"},
	  	"template": {"type": "string"},
      	"modules": {
		 	"type": "array",
		  	"items": {
			    "type": "object",
				"description": "Modules is a list of child modules to be rendered under this block",
			    "allOf": [
					{
						"properties": {
							"id": {"type": "string"}
						}
					},
					{
						"$ref": "#/definitions/module"
					}
		  	    ]
		  	}
	  	},
	  	"join": {"type": "string", "description": "Join is a template to use to join together modules.  Defaults to \" \"."}
	},
	"required": ["type"]
}`

// BlockModule renders a collection of other modules.
//
// Any module that outputs no text is considered "inactive" and will not be
// part of the result.
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

type blockModuleResult struct {
	// Modules is a map of results from executing each child module, indexed by
	// module ID.  Only modules that actually generated output will be included.
	Modules map[string]ModuleResult
	// ModuleArray is an array of results from executing each child module.  Only
	// modules that actually generated output will be included.
	ModuleArray []ModuleResult
}

// Execute the block module.
func (mod BlockModule) Execute(context *Context) ModuleResult {
	resultsArray := make([]ModuleResult, 0, len(mod.Modules))
	childDurations := perf.New(len(mod.Modules))
	resultsByID := make(map[string]ModuleResult, len(mod.Modules))

	for index := range mod.Modules {
		item := &mod.Modules[index]

		moduleDescription := fmt.Sprintf("%s(%d:%d)",
			item.ID,
			item.Line,
			item.Column,
		)

		childDurations.Start(moduleDescription)
		result := item.Module.Execute(context)
		childDurations.EndWithChildren(moduleDescription, result.ChildDurations)

		if len(result.Text) != 0 {
			resultsArray = append(resultsArray, result)
			if item.ID != "" {
				resultsByID[item.ID] = result
			}
		}
	}

	defaultText := mod.joinChildren(context, resultsArray)

	data := blockModuleResult{
		Modules:     resultsByID,
		ModuleArray: resultsArray,
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
	// Globals are the global variables.
	Globals *Globals
	// PrevColors is an `{FG, BG}` object containing color strings for the previous module's end style.
	PrevColors styling.CharacterColors
	// NextColors is an `{FG, BG}` object containing color strings for the next module's start style.
	NextColors styling.CharacterColors
	// Index is the index of the next module in the Modules array.
	Index int
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
			join, err = modtemplate.CompileTemplate(&context.Styles, context.Environment, "join", mod.Join)
			if err != nil {
				join = nil
			}
		}

		for index, child := range children {
			if join != nil && index != 0 {
				prev := children[index-1]
				joiner, err := modtemplate.TemplateToString(join, blockJoinData{
					Globals:    &context.Globals,
					PrevColors: prev.EndStyle,
					NextColors: child.StartStyle,
					Index:      index,
				})
				if err != nil {
					log.Warn(err.Error())
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
	registerModule(
		"block",
		registeredModule{
			jsonSchema: blockSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := BlockModule{Join: " "}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
