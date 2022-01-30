package modules

import (
	"strings"
	"text/template"

	"github.com/jwalton/kitsch/internal/kitsch/log"
	"github.com/jwalton/kitsch/internal/kitsch/modtemplate"
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"github.com/jwalton/kitsch/internal/kitsch/styling"
	"github.com/jwalton/kitsch/internal/perf"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas BlockModule

// BlockModule renders a collection of other modules.
//
// Any module that outputs no text is considered "inactive" and will not be
// part of the result.
//
type BlockModule struct {
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=block"`
	// Modules is a list of child modules to be rendered under this block
	Modules []ModuleWrapper `yaml:"modules" jsonschema:",required,ref=ModulesList"`
	// Join is a template to use to join together modules.  Defaults to " ".
	// This will be executed with template data of the form `{ PrevColors, NextColors, Index }`,
	// where PrevColors is the FG and BG color of last character of the previous
	// module, NextColors is the FG and BG color of the first character of the
	// next module, and Index is the index of the current module in the modules
	// array.
	Join string
}

type blockModuleResult struct {
	// Modules is a map of results from executing each child module, indexed by
	// module ID.  Only modules that actually generated output will be included.
	Modules map[string]ModuleWrapperResult
	// ModuleArray is an array of results from executing each child module.  Only
	// modules that actually generated output will be included.
	ModuleArray []ModuleWrapperResult
}

// Execute the block module.
func (mod BlockModule) Execute(context *Context) ModuleResult {
	resultsArray := make([]ModuleWrapperResult, 0, len(mod.Modules))
	childDurations := perf.New(len(mod.Modules))
	resultsByID := make(map[string]ModuleWrapperResult, len(mod.Modules))

	for index := range mod.Modules {
		wrapper := &mod.Modules[index]

		moduleDescription := wrapper.String()
		childDurations.Start(moduleDescription)
		result := wrapper.Execute(context)
		childDurations.EndWithChildren(moduleDescription, result.Performance)

		if len(result.Text) != 0 {
			resultsArray = append(resultsArray, result)
			if wrapper.config.ID != "" {
				resultsByID[wrapper.config.ID] = result
			}
		}
	}

	defaultText := mod.joinChildren(context, resultsArray)

	result := ModuleResult{
		DefaultText: defaultText,
		Performance: childDurations,
		Data: blockModuleResult{
			Modules:     resultsByID,
			ModuleArray: resultsArray,
		},
	}

	if len(resultsArray) > 0 {
		result.StartStyle = styling.CharacterColors{
			FG: defaultString(result.StartStyle.FG, resultsArray[0].StartStyle.FG),
			BG: defaultString(result.StartStyle.BG, resultsArray[0].StartStyle.BG),
		}
		lastChild := len(resultsArray) - 1
		result.EndStyle = styling.CharacterColors{
			FG: defaultString(result.EndStyle.FG, resultsArray[lastChild].EndStyle.FG),
			BG: defaultString(result.EndStyle.BG, resultsArray[lastChild].EndStyle.BG),
		}
	}

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

func (mod BlockModule) joinChildren(context *Context, children []ModuleWrapperResult) string {
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
			jsonSchema: schemas.BlockModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := BlockModule{Join: " "}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
