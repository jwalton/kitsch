// Package condition is a reusable "condition" object which can be used in configuration
// files to indicate when something should be done.
package condition

import (
	"runtime"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
)

// Conditions represents a condition which can be used in configuration files to
// specify when a module or project should be used.
type Conditions struct {
	// IfAncestorFiles is a list of files to search for in the project folder,
	// or another folder higher up in the directory structure.
	IfAncestorFiles []string `yaml:"ifAncestorFiles"`
	// IfFiles is a list of files to search for in the project folder.
	IfFiles []string `yaml:"ifFiles"`
	// IfExtensions is a list of extensions to search for in the project folder.
	IfExtensions []string `yaml:"ifExtensions"`
	// IfOS is a list of operating systems.  If the current GOOS is not in
	// the list, then this project type is not matched.
	IfOS []string `yaml:"ifOS"`
	// IfNotOS is a list of operating systems.  If the current GOOS is in
	// the list, then this project type is not matched.
	IfNotOS []string `yaml:"ifNotOS"`
}

// IsEmpty returns true if the condition has no conditions to match.
func (conditions Conditions) IsEmpty() bool {
	return len(conditions.IfAncestorFiles) == 0 &&
		len(conditions.IfFiles) == 0 &&
		len(conditions.IfExtensions) == 0 &&
		len(conditions.IfOS) == 0 &&
		len(conditions.IfNotOS) == 0
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
		if directory.HasFile(file) {
			return true
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
	if len(conditions.IfNotOS) > 0 {
		if contains(conditions.IfNotOS, runtime.GOOS) {
			return false
		}
	}

	if len(conditions.IfOS) > 0 {
		return contains(conditions.IfOS, runtime.GOOS)
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
