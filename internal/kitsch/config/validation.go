package config

import (
	"bytes"
	"encoding/json"
	"strings"

	// For baseSchema.
	_ "embed"
	"text/template"

	"github.com/jwalton/kitsch/internal/kitsch/condition"
	"github.com/jwalton/kitsch/internal/kitsch/getters"
	"github.com/jwalton/kitsch/internal/kitsch/modules"
	"github.com/jwalton/kitsch/internal/kitsch/projects"
	"github.com/jwalton/kitsch/internal/kitsch/schemautils"
)

//go:embed jsonschema.json
var schemaTemplate string

// JSONSchema returns the schema for the configuration file.
func JSONSchema() string {
	tmpl := template.Must(template.New("jsonSchema").Parse(schemaTemplate))

	defs := strings.Join([]string{
		getters.JSONSchemaDefinitions,
		condition.JSONSchemaDefinitions,
		projects.JSONSchemaDefinitions,
		modules.JSONSchemaDefinitions(),
	}, ",\n")

	data := map[string]string{
		"Definitions": defs,
	}

	var b bytes.Buffer
	err := tmpl.Execute(&b, data)
	if err != nil {
		panic(err)
	}
	rawSchema := b.String()

	// Parse and pretty-print rawSchema.
	var rawSchemaParsed map[string]interface{}
	err = json.Unmarshal([]byte(rawSchema), &rawSchemaParsed)
	if err != nil {
		return rawSchema
	}

	result, err := json.MarshalIndent(rawSchemaParsed, "", "    ")
	if err != nil {
		return rawSchema
	}
	return string(result)
}

// ValidateConfiguration validates the configuration file.
func ValidateConfiguration(yamlData []byte) error {
	// First try to load the configuration file.
	var config = Config{}
	err := config.LoadFromYaml(yamlData, true)
	if err != nil {
		return err
	}

	// TODO: Do custom validation for each module - validate all styles,
	// execute all templates, check for unknown fields, etc..

	// Validate the configuration file against the JSON schema.  This
	// will catch any errors in the configuration file, but tends to have
	// not-so-pretty error messages.
	err = schemautils.ValidateYamlAgainstSchema(yamlData, JSONSchema())
	if err != nil {
		return err
	}

	return nil
}
