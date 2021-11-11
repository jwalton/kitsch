// Package getters contains a Getter, which is an object that retrieves a value from
// the environment or file system.
package getters

import "github.com/jwalton/kitsch-prompt/internal/fileutils"

// Getter retrieves a text value from the file system or environment.
type Getter interface {
	// GetValue gets the value for this getter.  The return value will be either a string,
	// of if the value is a JSON, YAML, or TOML object, and the `ValueTemplate` is not set,
	// the parsed contents of the object.
	GetValue(folder fileutils.Directory) (interface{}, error)
}
