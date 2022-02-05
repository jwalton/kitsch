// Package projects is used to detect the project type of a directory.
//
package projects

import (
	"github.com/jwalton/kitsch/internal/kitsch/condition"
)

//go:generate go run ../genSchema/main.go --private ProjectType

// ProjectType represents configuration for checking if a folder is of a specific
// project type.
type ProjectType struct {
	// Name is the name of this project type.
	Name string `yaml:"name"`
	// Style is a default style for this project type.
	Style string `yaml:"style"`
	// Conditions are the conditions that must be met for this project type to be used.
	Conditions *condition.Conditions `yaml:"conditions,omitempty" jsonschema:",ref"`
	// ToolSymbol is the default symbol to use for this project type.
	ToolSymbol string `yaml:"toolSymbol"`
	// ToolVersion is used to retrieve the version of the build tool for this project.
	ToolVersion getterList `yaml:"toolVersion" jsonschema:",ref=GetterList"`
	// PackageManagerSymbol is the optional default symbol to use for the
	// package manager for this project type.
	PackageManagerSymbol string `yaml:"packageManagerSymbol"`
	// PackageManagerVersion is, if specified, used to retrieve the version of the
	// package manager for this project.
	PackageManagerVersion getterList `yaml:"packageManagerVersion"  jsonschema:",ref=GetterList"`
	// PackageVersion is, if specified, used to retrieve the version of the
	// project's package.
	PackageVersion getterList `yaml:"packageVersion" jsonschema:",ref=GetterList"`
}
