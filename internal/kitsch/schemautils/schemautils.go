package schemautils

import (
	"errors"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
)

// ValidateYamlNodeAgainstSchema validates a YAML node against a JSON schema.
func ValidateYamlNodeAgainstSchema(node *yaml.Node, schemaText string) error {
	var m interface{}
	err := node.Decode(&m)
	if err != nil {
		return err
	}
	return validateYamlAgainstSchema(m, schemaText)
}

// ValidateYamlAgainstSchema validates some YAML data against a JSON schema.
func ValidateYamlAgainstSchema(yamlData []byte, schemaText string) error {
	var m interface{}
	err := yaml.Unmarshal(yamlData, &m)
	if err != nil {
		return err
	}
	return validateYamlAgainstSchema(m, schemaText)
}

func validateYamlAgainstSchema(unmarshalledYaml interface{}, schemaText string) error {
	var err error
	unmarshalledYaml, err = toStringKeys(unmarshalledYaml)
	if err != nil {
		return err
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(schemaText)); err != nil {
		return err
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return err
	}

	if err := schema.Validate(unmarshalledYaml); err != nil {
		return err
	}

	return nil

}

func toStringKeys(val interface{}) (interface{}, error) {
	var err error
	switch val := val.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			k, ok := k.(string)
			if !ok {
				return nil, errors.New("found non-string key")
			}
			m[k], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case []interface{}:
		var l = make([]interface{}, len(val))
		for i, v := range val {
			l[i], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return l, nil
	default:
		return val, nil
	}
}
