package projects

import (
	"fmt"

	"github.com/jwalton/kitsch/internal/kitsch/getters"
	"github.com/jwalton/kitsch/internal/kitsch/log"
	"gopkg.in/yaml.v3"
)

type getterList []getters.Getter

// UnmarshalYAML unmarshals a YAML node into a ProjectTypeList.
func (list *getterList) UnmarshalYAML(node *yaml.Node) error {
	if node == nil {
		return fmt.Errorf("no value provided")
	}

	if node.Kind == yaml.MappingNode {
		// A single getter.
		getter := getters.CustomGetter{}
		err := node.Decode(&getter)
		if err != nil {
			return err
		}
		*list = []getters.Getter{getter}
		return nil
	}

	// A nodeList of getters.
	var nodeList []yaml.Node
	if err := node.Decode(&nodeList); err != nil {
		return err
	}

	newList := make([]getters.Getter, len(nodeList))
	for index := range nodeList {
		getter := getters.CustomGetter{}
		err := nodeList[index].Decode(&getter)
		if err != nil {
			return err
		}
		newList[index] = getter
	}
	*list = newList

	return nil
}

// MergeProjectTypes merges two sets of ProjectTypes.  Any ProjectTypes in the
// "to" set will be merged with the ProjectType with the same name in the
// "from" set, and if "addMissing" is true then any projects in the "from" set
// that aren't in the "to" set will be added to the end of the "to" set.
func MergeProjectTypes(to []ProjectType, from []ProjectType, addMissing bool) []ProjectType {
	usedMap := map[string]interface{}{}

	fromMap := map[string]*ProjectType{}
	for index := range from {
		fromMap[from[index].Name] = &from[index]
	}

	result := make([]ProjectType, 0, len(to)+len(from))

	// Copy over items from the "to", merging in the appropriate item from "from" if there is one.
	for _, toItem := range to {
		if _, ok := usedMap[toItem.Name]; ok {
			log.Warn(fmt.Sprintf("duplicate project type: %s", toItem.Name))
			continue
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

	return result
}

func mergeProjectType(to ProjectType, from ProjectType) ProjectType {
	if to.Style == "" {
		to.Style = from.Style
	}
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
