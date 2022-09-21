package modules

import (
	"fmt"
	"text/template"
	"time"

	"github.com/jwalton/kitsch/internal/kitsch/log"
	"github.com/jwalton/kitsch/internal/kitsch/modtemplate"
	"github.com/jwalton/kitsch/internal/kitsch/styling"
	"github.com/jwalton/kitsch/internal/perf"
	"gopkg.in/yaml.v3"
)

// TOOD: Store filename in ModuleWrapper?

// ModuleWrapper represents an item within a list of modules.
type ModuleWrapper struct {
	// config is common configuration for this module.
	config CommonConfig
	// Module is the actual module.
	Module Module
	// Line is the line number of the module in the configuration file.
	Line int
	// Column is the column number of the module in the configuration file.
	Column int
	// YamlNode is the YAML node that this module was read from, or nil if this module
	// was not loaded from YAML.
	YamlNode *yaml.Node
}

// ModuleWrapperResult represents the output of a ModuleWrapper.
type ModuleWrapperResult struct {
	// Text contains the rendered output of the module, either the default text
	// generated by the module itself, or the output from the template if one
	// was specified.
	Text string
	// Data contains any template data generated by the module.
	Data interface{}
	// StartStyle contains the foreground and background colors of the first
	// character in Text.  Note that this is based on the declared style for the
	// module - if the style for the module says the string should be colored
	// blue, but a template is used to change the color of the first character
	// to red, this will still say it is blue.
	StartStyle styling.CharacterColors
	// EndStyle is similar to StartStyle, but contains the colors of the last
	// character in Text.
	EndStyle styling.CharacterColors
	// Duration is the time it took this module to execute.
	Duration time.Duration
	// Performance is an array of execution times for children of this module.
	Performance *perf.Performance
}

// UnmarshalYAML converts a YAML node into a ModuleWrapper.
func (wrapper *ModuleWrapper) UnmarshalYAML(node *yaml.Node) error {
	if node == nil {
		return fmt.Errorf("no value provided")
	}
	wrapper.YamlNode = node
	wrapper.Line = node.Line
	wrapper.Column = node.Column

	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a map at (%d:%d)", node.Line, node.Column)
	}

	// Read the common configuration for this module, so we can work out what type it is.
	config, err := getCommonConfig(node)
	if err != nil {
		return err
	}
	wrapper.config = config

	// Load the actual module from the factory.
	mod, ok := registeredModules[config.Type]
	if !ok {
		return fmt.Errorf("unknown type %s (%d:%d)", config.Type, node.Line, node.Column)
	}
	module, err := mod.factory(node)
	if err != nil {
		return err
	}
	wrapper.Module = module

	return nil
}

func (wrapper ModuleWrapper) String() string {
	name := wrapper.config.Type
	if wrapper.config.ID != "" {
		name = name + "#" + wrapper.config.ID
	}

	return fmt.Sprintf("%s(%d:%d)",
		name,
		wrapper.Line,
		wrapper.Column,
	)
}

// Execute executes this module.  This will run the underlying Module, and then
// apply styling and the template from the CommonConfig.
func (wrapper ModuleWrapper) Execute(context *Context) ModuleWrapperResult {
	if !wrapper.config.Conditions.IsEmpty() && !wrapper.config.Conditions.Matches(context.Directory) {
		// If the item has conditions, and they don't match, return an empty result.
		return ModuleWrapperResult{}
	}

	// If the module has no timeout, use the default timeout.
	timeout := time.Duration(wrapper.config.Timeout) * time.Millisecond
	if timeout == 0 && wrapper.config.Type != "block" {
		timeout = context.DefaultTimeout
	}

	start := time.Now()

	// Run the module in a goroutine, so we can time it out.
	ch := make(chan ModuleWrapperResult, 1)
	go func() {
		moduleResult := wrapper.Module.Execute(context)
		ch <- processModuleResult(context, wrapper, moduleResult)
	}()

	var result ModuleWrapperResult
	if timeout <= 0 {
		result = <-ch
	} else {
		// If the module doesn't execute in time, return an empty result.
		select {
		case result = <-ch:
		case <-time.After(timeout):
			// Module timed out!
			// TODO: Record a list of which modules timed out in the context,
			// so we can display a list of them in a warning.
			log.Warn("Module ", wrapper.String(), " timed out after ", timeout)
			result = ModuleWrapperResult{}
		}
	}

	result.Duration = time.Since(start)

	return result
}

// TemplateData is the common data structure passed to a template when it is executed.
type TemplateData struct {
	// Text is the default text produced by this module
	Text string
	// Data is the data for this template.
	Data interface{}
	// Global is the global data.
	Globals *Globals
}

func compileModuleTemplate(context *Context, tmpl string) (*template.Template, error) {
	return modtemplate.CompileTemplate(context.Styles, context.Environment, "module-template", tmpl)
}

// executeModule is called to execute a module.  This handles "common" stuff that
// all modules do, like calling templates.
func processModuleResult(
	context *Context,
	moduleWrapper ModuleWrapper,
	moduleResult ModuleResult,
) ModuleWrapperResult {
	styleStr := moduleWrapper.config.Style
	if moduleResult.StyleOverride != "" {
		styleStr = moduleResult.StyleOverride
	}
	style := context.GetStyle(styleStr)

	text := moduleResult.DefaultText
	startStyle := moduleResult.StartStyle
	endStyle := moduleResult.EndStyle

	if moduleWrapper.config.Template != "" {
		tmpl, err := compileModuleTemplate(context, moduleWrapper.config.Template)
		if err != nil {
			log.Warn(fmt.Sprintf("Error compiling template in %s: %v", moduleWrapper.String(), err))
		} else {
			templateData := TemplateData{
				Data:    moduleResult.Data,
				Globals: &context.Globals,
				Text:    moduleResult.DefaultText,
			}

			text, err = modtemplate.TemplateToString(tmpl, templateData)
			if err != nil {
				log.Warn(fmt.Sprintf(
					"Error executing template in %s:\n%s\n%v",
					moduleWrapper.String(),
					moduleWrapper.config.Template,
					err,
				))
				text = moduleResult.DefaultText
			}
		}
	}

	if style != nil && text != "" {
		text, startStyle, endStyle = style.ApplyGetColors(text)
	}

	return ModuleWrapperResult{
		Text:        text,
		Data:        moduleResult.Data,
		StartStyle:  startStyle,
		EndStyle:    endStyle,
		Performance: moduleResult.Performance,
	}
}

// RenderPrompt renders the top-level module in a prompt.
func RenderPrompt(context *Context, root ModuleWrapper) (ModuleWrapperResult, string) {
	result := root.Execute(context)
	return result, processFlexibleSpaces(context.Globals.TerminalWidth, result.Text, context.FlexibleSpaceReplacement)
}

// func testTemplate(context *Context, prefix string, template string, dataMap map[string]interface{}) {
// 	if template == "" {
// 		return
// 	}

// 	tmpl, err := compileModuleTemplate(context, template)
// 	if err != nil {
// 		log.Warn(fmt.Sprintf("%s: Error compiling template: %v", prefix, err))
// 	} else {
// 		for description, data := range dataMap {
// 			templateData := TemplateData{
// 				Data:    data,
// 				Globals: &context.Globals,
// 				Text:    "",
// 			}

// 			_, err = modtemplate.TemplateToString(tmpl, templateData)
// 			if err != nil {
// 				log.Warn(fmt.Sprintf("%s: Error executing template with %s: %v", prefix, description, err))
// 			}
// 		}
// 	}
// }
