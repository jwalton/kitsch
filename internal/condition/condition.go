// Package condition is a reusable "condition" object which can be used in configuration
// files to indicate when something should be done.
package condition

import (
	"runtime"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
)

// Condition represents a condition which can be used in configuration files to
// specify when a module or project should be used.
type Condition struct {
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
func (condition Condition) IsEmpty() bool {
	return len(condition.IfAncestorFiles) == 0 &&
		len(condition.IfFiles) == 0 &&
		len(condition.IfExtensions) == 0 &&
		len(condition.IfOS) == 0 &&
		len(condition.IfNotOS) == 0
}

// Matches returns true if this condition is matched in the given directory
// and for the current operating system.
func (condition Condition) Matches(directory fileutils.Directory) bool {
	if !condition.matchesOS() {
		return false
	}

	for _, extension := range condition.IfExtensions {
		if directory.HasExtension(extension) {
			return true
		}
	}

	for _, file := range condition.IfFiles {
		if directory.HasFile(file) {
			return true
		}
	}

	for _, ancestorFile := range condition.IfAncestorFiles {
		result := directory.FindFileInAncestors(ancestorFile)
		if result != "" {
			return true
		}
	}

	return false
}

func (condition Condition) matchesOS() bool {
	if len(condition.IfNotOS) > 0 {
		if contains(condition.IfNotOS, runtime.GOOS) {
			return false
		}
	}

	if len(condition.IfOS) > 0 {
		return contains(condition.IfOS, runtime.GOOS)
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
