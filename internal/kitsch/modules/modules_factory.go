package modules

import (
	"fmt"
	"sort"
	"strings"

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
    "allOf": [
      { "$ref": "#/definitions/CommonConfig" },
      { "oneOf": [` + strings.Join(moduleRefs, ", ") + `] }
    ]
}`
	definitions = append(definitions, moduleDefinition)

	return strings.Join(definitions, ",\n")
}
