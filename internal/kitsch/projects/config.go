package projects

import (
	"fmt"

	"github.com/jwalton/kitsch-prompt/internal/kitsch/condition"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/getters"
	"gopkg.in/yaml.v3"
)

type projectTypeSpec struct {
	// Name is the name of this project type.
	Name string `yaml:"name"`
	// Condition is the condition that must be met for this project type to be used.
	Conditions condition.Conditions `yaml:"conditions"`
	// ToolSymbol is the default symbol to use for this project type.
	ToolSymbol string `yaml:"toolSymbol"`
	// PackageManagerSymbol is the default symbol to use for the package manager
	// for this project type.
	PackageManagerSymbol string `yaml:"packageManagerSymbol"`
	// ToolVersion is used to retrieve the version of the build tool for this project.
	ToolVersion getters.CustomGetter `yaml:"toolVersion"`
	// PackageManagerVersion is, if specified, used to retrieve the version of the
	// package manager for this project.
	PackageManagerVersion getters.CustomGetter `yaml:"packageManagerVersion"`
	// PackageVersion is, if specified, used to retrieve the version of the
	// project's package.
	PackageVersion getters.CustomGetter `yaml:"packageVersion"`
}

// UnmarshalYAML unmarshals a YAML node into a ProjectType.
func (item *ProjectType) UnmarshalYAML(node *yaml.Node) error {
	if node == nil {
		return fmt.Errorf("no value provided")
	}

	spec := projectTypeSpec{}
	err := node.Decode(&spec)
	if err != nil {
		return err
	}

	item.Name = spec.Name
	item.Conditions = spec.Conditions
	item.ToolSymbol = spec.ToolSymbol
	item.PackageManagerSymbol = spec.PackageManagerSymbol
	if spec.ToolVersion.Type != "" {
		item.ToolVersion = spec.ToolVersion
	}
	if spec.PackageManagerVersion.Type != "" {
		item.PackageManagerVersion = spec.PackageManagerVersion
	}
	if spec.PackageVersion.Type != "" {
		item.PackageVersion = spec.PackageVersion
	}

	return nil
}

// MergeProjectTypes merges two sets of ProjectTypes.  Any ProjectTypes in the
// "to" set will be merged with the ProjectType with the same name in the
// "from" set, and if "addMissing" is true then any projects in the "from" set
// that aren't in the "to" set will be added to the end of the "to" set.
func MergeProjectTypes(to []ProjectType, from []ProjectType, addMissing bool) ([]ProjectType, error) {
	usedMap := map[string]interface{}{}

	fromMap := map[string]*ProjectType{}
	for index := range from {
		fromMap[from[index].Name] = &from[index]
	}

	result := make([]ProjectType, 0, len(to)+len(from))

	// Copy over items from the "to", merging in the appropriate item from "from" if there is one.
	for _, toItem := range to {
		if _, ok := usedMap[toItem.Name]; ok {
			return nil, fmt.Errorf("duplicate project type: %s", toItem.Name)
		}
		usedMap[toItem.Name] = nil

		if fromItem, ok := fromMap[toItem.Name]; ok {
			toItem = mergeProjectType(toItem, *fromItem)
		}
		result = append(result, toItem)
	}

	// Add in any missing items from "from".
	if addMissing {
		for _, fromItem := range from {
			if _, ok := usedMap[fromItem.Name]; !ok {
				result = append(result, fromItem)
			}
		}
	}

	return result, nil
}

func mergeProjectType(to ProjectType, from ProjectType) ProjectType {
	if to.Conditions.IsEmpty() {
		to.Conditions = from.Conditions
	}
	if to.ToolSymbol == "" {
		to.ToolSymbol = from.ToolSymbol
	}
	if to.ToolVersion == nil {
		to.ToolVersion = from.ToolVersion
	}
	if to.PackageManagerSymbol == "" {
		to.PackageManagerSymbol = from.PackageManagerSymbol
	}
	// TODO: How do we remove an optional value?  Maybe we could put `{}` in the YAML,
	// and load a "NullGetter" that always returns ""?
	if to.PackageManagerVersion == nil {
		to.PackageManagerVersion = from.PackageManagerVersion
	}
	if to.PackageVersion == nil {
		to.PackageVersion = from.PackageVersion
	}

	return to
}
