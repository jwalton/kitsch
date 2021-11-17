package env

import "github.com/jwalton/kitsch-prompt/internal/gitutils"

// DummyEnv is a dummy environment for use in unit testing.
type DummyEnv struct {
	// Env contains the environment variables for this dummy environment.
	Env map[string]string
	// Root is true if this
	Root bool
	// TestJobs is the return value for Jobs().
	TestJobs int
	// TestGit is the return value for Git().
	TestGit *gitutils.GitUtils
}

// Getenv returns the value of the specifed environment variable.
// Returns the value set in `env.Env`.
func (env *DummyEnv) Getenv(key string) string {
	val, ok := env.Env[key]
	if ok {
		return val
	}
	return ""
}

// IsRoot returns true if DummyEnv.Root is set.
func (env *DummyEnv) IsRoot() bool {
	return env.Root
}

// HasSomeEnv returns true if one or more of the given environment variables are set.
func (env *DummyEnv) HasSomeEnv(keys ...string) bool {
	for _, key := range keys {
		if env.Getenv(key) != "" {
			return true
		}
	}
	return false
}

// Jobs returns env.TestJobs.
func (env *DummyEnv) Jobs() int {
	return env.TestJobs
}

// Git returns a git instance for the current repo, or nil if the current
// working directory is not part of a git repo, or git is uninstalled.
func (env *DummyEnv) Git() *gitutils.GitUtils {
	return env.TestGit
}
