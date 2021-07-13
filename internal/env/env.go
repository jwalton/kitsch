// Package env provides a thread safe interface between modules and the runtime
// environment, which can be easily swapped out for unit test purposes.
package env

import (
	"os"
	"os/user"
	"runtime"
)

// Env is an interface between modules and the runtime environment.
type Env interface {
	// Getenv returns the value of the specifed environment variable.
	Getenv(key string) string
	// GetUsername returns the current user's username.  Returns the empty string
	// if the username cannot be determined.
	GetUsername() string
	// IsRoot returns true if the current user is root.
	IsRoot() bool
	// HasSomeEnv returns true if at least one of the specified environment variables
	// is defined and non-empty.  For example:
	//
	//     env.HasSomeEnv("SSH_CLIENT", "SSH_CONNECTION", "SSH_TTY")
	//
	// would return true if this is an SSH session.
	HasSomeEnv(...string) bool
	// Getwd returns the current working directory.
	Getwd() string
	// UserHomeDir returns the current user's home directory.
	UserHomeDir() string
	// Jobs returns the number of jobs running in the background.
	Jobs() int
	// CmdDuration returns the duration of the last run command, in milliseconds.
	CmdDuration() int
	// Status returns the exit code of the last command run.
	Status() int
	// Keymap returns the zsh/fish keymap.
	Keymap() string
}

type defaultEnv struct {
	jobs        int
	cmdDuration int
	status      int
	keymap      string
}

// NewEnv creates a new instance of Env.
func NewEnv(
	jobs int,
	cmdDuration int,
	status int,
	keymap string,
) Env {
	return &defaultEnv{
		jobs:        jobs,
		cmdDuration: cmdDuration,
		status:      status,
		keymap:      keymap,
	}
}

func (*defaultEnv) Getenv(key string) string {
	return os.Getenv(key)
}

func (*defaultEnv) GetUsername() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Name
}

func (*defaultEnv) IsRoot() bool {
	// TODO: How to handle plan9 here?
	if runtime.GOOS == "windows" {
		return false
	}

	user, err := user.Current()
	if err != nil {
		return false
	}
	return user.Uid == "0"
}

func (*defaultEnv) HasSomeEnv(keys ...string) bool {
	for _, key := range keys {
		if os.Getenv(key) != "" {
			return true
		}
	}
	return false
}

func (*defaultEnv) Getwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func (*defaultEnv) UserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~"
	}
	return home
}

func (env *defaultEnv) Jobs() int {
	return env.jobs
}

func (env *defaultEnv) CmdDuration() int {
	return env.cmdDuration
}

// Status returns the exit code of the last command run.
func (env *defaultEnv) Status() int {
	return env.status
}

// Keymap returns the zsh/fish keymap
func (env *defaultEnv) Keymap() string {
	return env.keymap
}
