// Package condition is a reusable "condition" object which can be used in configuration
// files to indicate when something should be done.
package condition

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jwalton/kitsch/internal/fileutils"
)

//go:generate go run ../genSchema/main.go --private Conditions

// Conditions represents a condition which can be used in configuration files to
// specify when a module or project should be used.
type Conditions struct {
	// IfAncestorFiles is a list of files to search for in the current folder,
	// or another folder higher up in the directory structure.
	IfAncestorFiles []string `yaml:"ifAncestorFiles"`
	// IfFiles is a list of files to search for in the current folder.
	IfFiles []string `yaml:"ifFiles"`
	// IfExtensions is a list of extensions to search for in the current folder.
	IfExtensions []string `yaml:"ifExtensions"`
	// OnlyIfOS is a list of operating systems.  If the current GOOS is not in
	// the list, then the Conditions are not met, even if other conditions would
	// be satisfied.
	OnlyIfOS []string `yaml:"onlyIfOS"`
	// OnlyIfNotOS is a list of operating systems.  If the current GOOS is in
	// the list, then the Conditions are not met, even if other conditions would
	// be satisfied.
	OnlyIfNotOS []string `yaml:"onlyIfNotOS"`
}

// IsEmpty returns true if the condition has no conditions to match.
func (conditions Conditions) IsEmpty() bool {
	return len(conditions.IfAncestorFiles) == 0 &&
		len(conditions.IfFiles) == 0 &&
		len(conditions.IfExtensions) == 0 &&
		len(conditions.OnlyIfOS) == 0 &&
		len(conditions.OnlyIfNotOS) == 0
}

// Matches returns true if this condition is matched in the given directory
// and for the current operating system.
func (conditions Conditions) Matches(directory fileutils.Directory) bool {
	if !conditions.matchesOS() {
		return false
	}

	for _, extension := range conditions.IfExtensions {
		if directory.HasExtension(extension) {
			return true
		}
	}

	for _, file := range conditions.IfFiles {
		if filepath.IsAbs(file) || strings.HasPrefix(file, "..") {
			// If the file is an absolute path or a parent directory,
			// we need to go directly to the OS to see if it exists.
			if fileutils.FileExists(filepath.Join(directory.Path(), file)) {
				return true
			}
		} else {
			if directory.HasFile(file) {
				return true
			}
		}
	}

	for _, ancestorFile := range conditions.IfAncestorFiles {
		result := directory.FindFileInAncestors(ancestorFile)
		if result != "" {
			return true
		}
	}

	return false
}

func (conditions Conditions) matchesOS() bool {
	if len(conditions.OnlyIfNotOS) > 0 {
		if contains(conditions.OnlyIfNotOS, runtime.GOOS) {
			return false
		}
	}

	if len(conditions.OnlyIfOS) > 0 {
		return contains(conditions.OnlyIfOS, runtime.GOOS)
	}

	return true
}

// TODO: Replace this with a generic implementation.
func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

//JSONSchemaDefinitions is a string containing JSON schema definitions for objects in the conditions package.
var JSONSchemaDefinitions = "\"Conditions\": " + conditionsJSONSchema
