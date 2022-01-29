package gen

import (
	"strings"

	"github.com/fatih/structtag"
)

type schemaTags struct {
	required  bool
	ref       bool
	refStruct string
	enum      []string
	fieldName string
}

func parseStructTags(fieldName string, fieldTags string) schemaTags {
	// Strip the `s from the tag value.
	tags, err := structtag.Parse(fieldTags)
	if err != nil {
		panic(err.Error() + ": " + fieldTags)
	}

	fieldName = getJSONSchemaFieldName(fieldName, tags)
	result := schemaTags{fieldName: fieldName}

	for _, tag := range tags.Tags() {
		if tag.Key == "jsonschema" {
			for _, option := range tag.Options {
				if option == "required" {
					result.required = true
				} else if option == "ref" {
					result.ref = true
				} else if strings.HasPrefix(option, "ref=") {
					result.ref = true
					result.refStruct = option[4:]
				} else if strings.HasPrefix(option, "enum=") {
					result.enum = strings.Split(option[5:], ":")
				} else {
					panic("Unknown jsonschema tag: " + option)
				}
			}
		}
	}

	return result
}

func getJSONSchemaFieldName(fieldName string, tags *structtag.Tags) string {
	if tag, err := tags.Get("jsonschema"); err == nil && tag.Name != "" {
		return tag.Name
	}

	if tag, err := tags.Get("json"); err == nil {
		return tag.Name
	}

	if tag, err := tags.Get("yaml"); err == nil {
		return tag.Name
	}

	return strings.ToLower(fieldName[0:1]) + fieldName[1:]
}
