// Package env provides a thread safe interface between modules and the runtime
// environment, which can be easily swapped out for unit test purposes.
package env

import (
	"os"
	"os/user"
	"runtime"

	"github.com/jwalton/kitsch-prompt/internal/gitutils"
)

// Env is an interface between modules and the runtime environment.
type Env interface {
	// Getenv returns the value of the specifed environment variable.
	Getenv(key string) string
	// IsRoot returns true if the current user is root.
	IsRoot() bool
	// HasSomeEnv returns true if at least one of the specified environment variables
	// is defined and non-empty.  For example:
	//
	//     env.HasSomeEnv("SSH_CLIENT", "SSH_CONNECTION", "SSH_TTY")
	//
	// would return true if this is an SSH session.
	HasSomeEnv(...string) bool
	// Jobs returns the number of jobs running in the background.
	Jobs() int
	// Git returns a git instance for the current repo, or nil if the current
	// working directory is not part of a git repo, or git is not installed.
	Git() *gitutils.GitUtils
}

type defaultEnv struct {
	cwd            string
	jobs           int
	gitInitialized bool
	git            *gitutils.GitUtils
}

// Options is a struct that can be used to customize the behavior of Env.
type Options struct {
	// CWD is the current working directory.
	CWD string
	// Jobs is the number of jobs running in the background.
	Jobs int
}

// New creates a new instance of Env.
func New(options Options) Env {
	return &defaultEnv{
		cwd:  options.CWD,
		jobs: options.Jobs,
	}
}

func (*defaultEnv) Getenv(key string) string {
	return os.Getenv(key)
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

func (env *defaultEnv) Jobs() int {
	return env.jobs
}

// Git returns a git instance for the current repo, or nil if the current
// working directory is not part of a git repo, or git is uninstalled.
func (env *defaultEnv) Git() *gitutils.GitUtils {
	if !env.gitInitialized {
		env.git = gitutils.New("git", env.cwd)
		env.gitInitialized = true
	}
	return env.git
}
