// Package env provides an interface to environment variables.
package env

import (
	"os"
)

// Env is an interface to environment variables which can be overridden for testing.
type Env interface {
	// Getenv returns the value of the specified environment variable, or an empty
	// string if the variable does not exist.
	Getenv(key string) string
	// HasSomeEnv returns true if at least one of the specified environment variables
	// is defined and non-empty.  For example:
	//
	//     env.HasSomeEnv("SSH_CLIENT", "SSH_CONNECTION", "SSH_TTY")
	//
	// would return true if this is an SSH session.
	HasSomeEnv(...string) bool
}

type defaultEnv struct{}

// New creates a new instance of Env.
func New() Env {
	return defaultEnv{}
}

func (defaultEnv) Getenv(key string) string {
	return os.Getenv(key)
}

func (defaultEnv) HasSomeEnv(keys ...string) bool {
	for _, key := range keys {
		if os.Getenv(key) != "" {
			return true
		}
	}
	return false
}
