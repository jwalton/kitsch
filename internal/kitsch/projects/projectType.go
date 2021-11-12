// Package projects is used to detect the project type of a directory.
//
package projects

import (
	"github.com/jwalton/kitsch-prompt/internal/kitsch/condition"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/getters"
)

// ProjectType represents configuration for checking if a folder is of a specific
// project type.
type ProjectType struct {
	// Name is the name of this project type.
	Name string
	// Conditions are the conditions that must be met for this project type to be used.
	Conditions condition.Conditions
	// ToolSymbol is the default symbol to use for this project type.
	ToolSymbol string
	// ToolVersion is used to retrieve the version of the build tool for this project.
	ToolVersion getters.Getter
	// PackageManagerSymbol is the optional default symbol to use for the
	// package manager for this project type.
	PackageManagerSymbol string
	// PackageManagerVersion is, if specified, used to retrieve the version of the
	// package manager for this project.
	PackageManagerVersion getters.Getter
	// PackageVersion is, if specified, used to retrieve the version of the
	// project's package.
	PackageVersion getters.Getter
}
