package env

// DummyEnv is a dummy environment for use in unit testing.
type DummyEnv struct {
	// Env contains the environment variables for this dummy environment.
	Env map[string]string
	// Root is true if this
	Root bool
	// CWD is the current working directory.
	CWD string
	// TestJobs is the return value for Jobs().
	TestJobs int
	// TestCmdDuration is the return value for CmdDuration().
	TestCmdDuration int
	// TestStatus is the return value for Status().
	TestStatus int
	// TestKeymap is the return value for Keymap().
	TestKeymap string
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

// GetUsername returns the value of the "USER" environment variable.
func (env *DummyEnv) GetUsername() string {
	return env.Getenv("USER")
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

// Getwd returns env.CWD, or "/Users/jwalton" if it is unset.
func (env *DummyEnv) Getwd() string {
	if env.CWD == "" {
		return "/Users/jwalton"
	}
	return env.CWD
}

// UserHomeDir returns the contents of the "HOME" environment variable.
func (env *DummyEnv) UserHomeDir() string {
	return env.Getenv("HOME")
}

// Jobs returns env.TestJobs.
func (env *DummyEnv) Jobs() int {
	return env.TestJobs
}

// CmdDuration returns the duration of the last run command, in milliseconds.
// Returns env.TestCmdDuration.
func (env *DummyEnv) CmdDuration() int {
	return env.TestCmdDuration
}

// Status returns the exit code of the last command run.
func (env *DummyEnv) Status() int {
	return env.TestStatus
}

// Keymap returns the zsh/fish keymap
func (env *DummyEnv) Keymap() string {
	return env.TestKeymap
}
