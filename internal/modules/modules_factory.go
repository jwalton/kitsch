package modules

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type typedYamlNode struct {
	Type string `yaml:"type"`
	ID   string `yaml:"id"`
}

type moduleFactory func(node *yaml.Node) (Module, error)

// moduleFactories lists factories for converting YAML nodes into modules.
// Modules register a factory via registerFactor().
var moduleFactories = map[string]moduleFactory{}

// Regsiter a factory to create the specified type of module.
func registerFactory(name string, factory moduleFactory) {
	if _, ok := moduleFactories[name]; ok {
		panic("Duplicate module factory registration: " + name)
	}
	moduleFactories[name] = factory
}

// ModuleSpec represents an item within a list of modules.
type ModuleSpec struct {
	// ID is a unique ID for this module within a ModuleList, if provided.
	ID string
	// Module is the actual module.
	Module Module
}

// UnmarshalYAML unmarshals a YAML node into a module.
func (item *ModuleSpec) UnmarshalYAML(node *yaml.Node) error {
	if node == nil {
		return fmt.Errorf("no value provided")
	}

	// Special case where node is a bare string.
	if node.Kind == yaml.ScalarNode && node.Tag == "!!str" {
		item.Module = TextModule{
			Text: node.Value,
		}
		return nil
	}

	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a map at (%d:%d)", node.Line, node.Column)
	}

	// Figure out the type of this node.
	moduleType, id, err := getTypeAndID(node)
	if err != nil {
		return err
	}

	factory := moduleFactories[moduleType]
	if factory == nil {
		return fmt.Errorf("unknown type %s (%d:%d)", moduleType, node.Line, node.Column)
	}

	module, err := factory(node)
	if err != nil {
		return err
	}

	if id != "" {
		item.ID = id
	} else {
		item.ID = moduleType
	}

	item.Module = module
	return nil
}

// getTypeAndID retrieves the "type" and "id" key of a YAML mapping node.
func getTypeAndID(node *yaml.Node) (string, string, error) {
	if node == nil {
		return "", "", fmt.Errorf("cannot get type of empty node")
	}

	var t typedYamlNode
	err := node.Decode(&t)

	if t.Type == "" {
		return "", "", fmt.Errorf("object is missing type (%d:%d)", node.Line, node.Column)
	}

	return t.Type, t.ID, err
}
