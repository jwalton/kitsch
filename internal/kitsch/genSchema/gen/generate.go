package gen

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

var defaultIndent = "    "

// TODO: Add min/max restrictions on the various int types? Or allow these to be
// configured through struct tags?
var basicTypes = map[string]string{
	"bool":    `boolean`,
	"string":  `string`,
	"int":     `integer`,
	"byte":    `integer`,
	"int8":    `integer`,
	"int16":   `integer`,
	"int32":   `integer`,
	"int64":   `integer`,
	"uint":    `integer`,
	"uint8":   `integer`,
	"uint16":  `integer`,
	"uint32":  `integer`,
	"uint64":  `integer`,
	"rune":    `integer`,
	"float32": `number`,
	"float64": `number`,
}

type schemaBuilder struct {
	pkg                  *packages.Package
	properties           []string
	required             []string
	indent               string
	additionalProperties bool
}

func newSchemaBuilder(filename string, additionalProperties bool) (*schemaBuilder, error) {
	cfg := &packages.Config{Mode: packages.NeedName |
		packages.NeedFiles |
		packages.NeedSyntax |
		packages.NeedTypes |
		packages.NeedTypesInfo}
	pkgs, err := packages.Load(cfg, "file="+filename)

	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("Could not read package for %s", filename)
	}

	return &schemaBuilder{
		pkg:                  pkgs[0],
		indent:               "",
		additionalProperties: additionalProperties,
	}, nil
}

// GenerateSchemaForStruct generates the JSON schema for a structure in a file.
func GenerateSchemaForStruct(filename string, structName string, additionalProperties bool) (string, error) {
	builder, err := newSchemaBuilder(filename, additionalProperties)
	if err != nil {
		return "", err
	}

	obj := builder.pkg.Types.Scope().Lookup(structName)
	struc, ok := obj.(*types.TypeName).Type().Underlying().(*types.Struct)
	if !ok || struc == nil {
		return "", fmt.Errorf("%s is not a struct", structName)
	}

	err = builder.addStruct(structName, struc)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (builder *schemaBuilder) String() string {
	indent := builder.indent

	lf := "\n" + indent

	parts := []string{
		`  "type": "object"`,
		`  "properties": {` + lf + `    ` + strings.Join(builder.properties, ","+lf+`    `) + lf + `  }`,
	}
	if len(builder.required) > 0 {
		parts = append(parts, `  "required": [`+quotedStrings(builder.required)+`]`)
	}
	if !builder.additionalProperties {
		parts = append(parts, `  "additionalProperties": false`)
	}

	result := indent + "{" + lf
	result += strings.Join(parts, `,`+lf)
	result += `}`

	return result
}

// addStruct adds all fields from a struct to the given schema builder.
func (builder *schemaBuilder) addStruct(structName string, s *types.Struct) error {
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)

		// Handle struct tag options
		tags := parseStructTags(field.Name(), s.Tag(i))

		// Skip fields which are not exported, or which should be skipped.
		if !field.Exported() || tags.fieldName == "-" {
			continue
		}

		// Add all fields from embedded structs directly to this struct.
		// We could instead use `allOf` here to merge the schemas, but then
		// we can't use "additionFields: false".
		if field.Embedded() {
			err := builder.addStruct("", field.Type().Underlying().(*types.Struct))
			if err != nil {
				return err
			}
			continue
		}

		var fieldSchema string
		if tags.ref {
			refStruct := tags.refStruct
			if refStruct == "" {
				var err error
				refStruct, err = getBareTypeName(field.Type())
				if err != nil {
					return err
				}
			}
			fieldSchema = fmt.Sprintf(`{"$ref": "#/definitions/%s"}`, refStruct)
		} else {
			description := builder.getDescriptionForField(structName, field.Name())
			description = strings.TrimSuffix(description, "\n")
			description = strings.ReplaceAll(description, "\n", " ")
			description = strings.ReplaceAll(description, `"`, `\"`)

			var err error
			fieldSchema, err = builder.generateSchemaForType(field.Type(), tags, description)
			if err != nil {
				return err
			}
		}

		// Add the property to the schema.
		builder.properties = append(builder.properties, fmt.Sprintf(`"%s": %s`, tags.fieldName, fieldSchema))
		if tags.required {
			builder.required = append(builder.required, tags.fieldName)
		}
	}

	return nil
}

// generateSchemaForType will generate a JSON schema for the specified type.
func (builder *schemaBuilder) generateSchemaForType(
	t types.Type,
	tags schemaTags,
	description string,
) (string, error) {
	var fieldSchema string

	// Generate the schema for the field
	if tags.enum != nil {
		fieldSchema = generateSchemaForEnum(tags.enum, description)
	} else {
		switch v := t.Underlying().(type) {
		case *types.Basic:
			basicSchemaType, ok := basicTypes[v.String()]
			if !ok {
				return "", fmt.Errorf("Unhandled basic type %s", v.String())
			}
			fieldSchema = fmt.Sprintf(`{"type": "%s", "description": "%s"}`, basicSchemaType, description)
		case *types.Struct:
			childBuilder := schemaBuilder{pkg: builder.pkg, indent: builder.indent + defaultIndent}
			bareTypeName, err := getBareTypeName(t)
			if err != nil {
				return "", err
			}
			err = childBuilder.addStruct(bareTypeName, v)
			if err != nil {
				return "", err
			}
			fieldSchema = childBuilder.String()
		case *types.Map:
			valueType := v.Elem()
			valueSchema, err := builder.generateSchemaForType(valueType, schemaTags{}, "")
			if err != nil {
				return "", err
			}
			fieldSchema = fmt.Sprintf(`{"type": "object", "description": "%s", "additionalProperties": %s}`, description, valueSchema)
		case *types.Array:
			valueType := v.Elem()
			valueSchema, err := builder.generateSchemaForType(valueType, schemaTags{}, "")
			if err != nil {
				return "", err
			}
			fieldSchema = fmt.Sprintf(`{"type": "array", "description": "%s", "items": %s}`, description, valueSchema)
		case *types.Slice:
			valueType := v.Elem()
			valueSchema, err := builder.generateSchemaForType(valueType, schemaTags{}, "")
			if err != nil {
				return "", err
			}
			fieldSchema = fmt.Sprintf(`{"type": "array", "description": "%s", "items": %s}`, description, valueSchema)
		default:
			return "", fmt.Errorf("Unhandled type: %T", t)
		}
	}

	return fieldSchema, nil
}

func generateSchemaForEnum(enum []string, description string) string {
	return fmt.Sprintf(`{"type": "string", "description": "%s", "enum": [%s]}`, description, quotedStrings(enum))
}

func quotedStrings(s []string) string {
	return `"` + strings.Join(s, `", "`) + `"`
}

func (builder *schemaBuilder) getDescriptionForField(
	structName string,
	fieldName string,
) string {
	// TODO: Must be a better way to do this.
	if structName == "" {
		return ""
	}

	for _, tree := range builder.pkg.Syntax {
		s := findStruct(tree, structName)
		if s != nil {
			for _, field := range s.Fields.List {
				if len(field.Names) > 0 && field.Names[0].Name == fieldName {
					return field.Doc.Text()
				}
			}
		}
	}

	return ""
}

func getBareTypeName(t types.Type) (string, error) {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	if named, ok := t.(*types.Named); ok {
		return named.Obj().Name(), nil
	}

	return "", fmt.Errorf("Unhandled type: %T", t)

}
