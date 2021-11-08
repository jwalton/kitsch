package projects

import (
	"fmt"

	"github.com/jwalton/kitsch-prompt/internal/condition"
	"github.com/jwalton/kitsch-prompt/internal/getters"
	"gopkg.in/yaml.v3"
)

type projectTypeSpec struct {
	// Name is the name of this project type.
	Name string `yaml:"name"`
	// Condition is the condition that must be met for this project type to be used.
	Condition condition.Condition `yaml:"condition"`
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
	item.Condition = spec.Condition
	item.ToolSymbol = spec.ToolSymbol
	item.PackageManagerSymbol = spec.PackageManagerSymbol
	item.ToolVersion = spec.ToolVersion
	if spec.PackageManagerVersion.Type != "" {
		item.PackageManagerVersion = spec.PackageManagerVersion
	}
	if spec.PackageVersion.Type != "" {
		item.PackageVersion = spec.PackageVersion
	}

	return nil
}
