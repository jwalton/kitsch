// Package getters contains a Getter, which is an object that retrieves a value from
// the environment or file system.
package getters

import (
	"github.com/jwalton/kitsch/internal/cache"
	"github.com/jwalton/kitsch/internal/fileutils"
)

// GetterContext is an interface used by a Getter to retrieve information from
// the environment or file system.
type GetterContext interface {
	// GetWorkingDirectory returns the current working directory.
	GetWorkingDirectory() fileutils.Directory

	// GetHomeDirectoryPath returns the path to the user's home directory.
	GetHomeDirectoryPath() string

	// Getenv returns the value of the specified environment variable.
	Getenv(key string) string

	// GetValueCache returns the value cache.
	GetValueCache() cache.Cache
}

// Getter retrieves a text value from the file system or environment.
type Getter interface {
	// GetValue gets the value for this getter.  The return value will be either a string,
	// of if the value is a JSON, YAML, or TOML object, and the `ValueTemplate` is not set,
	// the parsed contents of the object.
	GetValue(getterContext GetterContext) (interface{}, error)
}
