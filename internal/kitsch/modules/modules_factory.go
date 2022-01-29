package modules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jwalton/kitsch/internal/kitsch/condition"
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

type registeredModule struct {
	factory    func(node *yaml.Node) (Module, error)
	jsonSchema string
}

// registeredModules lists information about each type of module.
// Modules register a factory via registerModule().
var registeredModules = map[string]registeredModule{}

// Regsiter a factory to create the specified type of module.
func registerModule(name string, mod registeredModule) {
	if mod.jsonSchema == "" {
		panic("Missing JSON schema for module " + name)
	}
	if mod.factory == nil {
		panic("Missing factory for module " + name)
	}
	if _, ok := registeredModules[name]; ok {
		panic("Duplicate module factory registration: " + name)
	}
	registeredModules[name] = mod
}

// TOOD: Store filename in ModuleSpec?

// ModuleSpec represents an item within a list of modules.
type ModuleSpec struct {
	// ID is a unique ID for this module (within a list of modules), if provided.
	ID string
	// Conditions is a set of conditions that must be met for this module to be executed.
	Conditions condition.Conditions
	// Type is the type of this module.
	Type string
	// Module is the actual module.
	Module Module
	// Line is the line number of the module in the configuration file.
	Line int
	// Column is the column number of the module in the configuration file.
	Column int
	// Children is an array of child modules for this module.
	Children []ModuleSpec
	// YamlNode is the YAML node that this module was read from, or nil if this module
	// was not loaded from YAML.
	YamlNode *yaml.Node
}

// moduleSpecData is a subset of ModuleSpec that we read from the yanl.Node.
type moduleSpecData struct {
	Conditions condition.Conditions `yaml:"conditions"`
	ID         string               `yaml:"id"`
	Type       string               `yaml:"type"`
}

// UnmarshalYAML unmarshals a YAML node into a module.
func (item *ModuleSpec) UnmarshalYAML(node *yaml.Node) error {
	if node == nil {
		return fmt.Errorf("no value provided")
	}
	item.YamlNode = node

	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a map at (%d:%d)", node.Line, node.Column)
	}

	// Figure out the type of this node.
	data, err := getModuleSpecCommon(node)
	if err != nil {
		return err
	}

	mod, ok := registeredModules[data.Type]
	if !ok {
		return fmt.Errorf("unknown type %s (%d:%d)", data.Type, node.Line, node.Column)
	}

	module, err := mod.factory(node)
	if err != nil {
		return err
	}

	id := data.ID
	if id != "" {
		item.ID = id
	} else {
		item.ID = data.Type
	}

	item.Type = data.Type
	item.Conditions = data.Conditions
	item.Line = node.Line
	item.Column = node.Column
	item.Module = module

	// TODO: More generic handling of children?
	if block, ok := module.(BlockModule); ok {
		item.Children = block.Modules
	}

	return nil
}

// Execute executes this module.
func (item ModuleSpec) Execute(context *Context) ModuleResult {
	if !item.Conditions.IsEmpty() && !item.Conditions.Matches(context.Directory) {
		// If the item has conditions, and they don't match, return an empty result.
		return ModuleResult{}
	}
	return item.Module.Execute(context)
}

// getModuleSpecCommon retrieves the "type" and "id" key of a YAML mapping node.
func getModuleSpecCommon(node *yaml.Node) (moduleSpecData, error) {
	var t moduleSpecData
	if node == nil {
		return t, fmt.Errorf("cannot get type of empty node")
	}

	err := node.Decode(&t)
	if err != nil {
		return t, err
	}

	if t.Type == "" {
		return t, fmt.Errorf("object is missing type (%d:%d)", node.Line, node.Column)
	}

	return t, nil
}

// JSONSchemaForModule returns the JSON schema for a module.
func JSONSchemaForModule(typeName string) string {
	mod, ok := registeredModules[typeName]
	if !ok {
		panic("Unknown module type: " + typeName)
	}
	return mod.jsonSchema
}

// JSONSchemaDefinitions returns a string cotaining definitions to add to the JSON schema for all modules.
func JSONSchemaDefinitions() string {
	var definitions []string
	var moduleRefs []string

	definitions = append(definitions, fmt.Sprintf("\"CommonConfig\": %s", schemas.CommonConfigJSONSchema))

	keys := make([]string, 0, len(registeredModules))
	for name := range registeredModules {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	for _, name := range keys {
		mod := registeredModules[name]
		definitions = append(definitions, fmt.Sprintf("\"%s\": %s", name, mod.jsonSchema))
		moduleRefs = append(moduleRefs, fmt.Sprintf("{ \"$ref\": \"#/definitions/%s\" }", name))
	}

	// Add a "module" definition, which can be any module.
	moduleDefinition := `"module": {
    "type": "object",
	"required": [ "type" ],
    "oneOf": [` + strings.Join(moduleRefs, ", ") + `]
}`
	definitions = append(definitions, moduleDefinition)

	return strings.Join(definitions, ",\n")
}
