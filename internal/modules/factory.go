package modules

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type typedYamlNode struct {
	// The type of this node.
	Type string `yaml:"type"`
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

// CreateModule creates a module from a YAML node.
func CreateModule(node *yaml.Node) (Module, error) {
	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected a map at (%d:%d)", node.Line, node.Column)
	}

	// Figure out the type of this node.
	moduleType, err := getTypeFromNode(node)
	if err != nil {
		return nil, err
	}

	factory := moduleFactories[moduleType]
	if factory == nil {
		return nil, fmt.Errorf("unknown type %s (%d:%d)", moduleType, node.Line, node.Column)
	}

	return factory(node)
}

// ModuleList represents a list of modules.
type ModuleList struct {
	// The list of modules.
	Modules []Module `yaml:",inline"`
}

// UnmarshalYAML unmarshals a YAML node into a list of modules.
func (modules *ModuleList) UnmarshalYAML(node *yaml.Node) error {
	if node == nil {
		return fmt.Errorf("no value provided")
	}

	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("expected a sequence at (%d:%d)", node.Line, node.Column)
	}

	for _, subnode := range node.Content {
		m, err := CreateModule(subnode)
		if err != nil {
			return err
		}

		modules.Modules = append(modules.Modules, m)
	}

	return nil
}

// getTypeFromNode retrieves the "type" key of a YAML mapping node.
func getTypeFromNode(node *yaml.Node) (string, error) {
	if node == nil {
		return "", fmt.Errorf("cannot get type of empty node")
	}

	var t typedYamlNode
	err := node.Decode(&t)

	if t.Type == "" {
		return "", fmt.Errorf("object is missing type (%d:%d)", node.Line, node.Column)
	}

	return t.Type, err
}
