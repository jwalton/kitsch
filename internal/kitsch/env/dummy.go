package env

// DummyEnv is a dummy environment for use in unit testing.
type DummyEnv struct {
	// Env contains the environment variables for this dummy environment.
	Env map[string]string
}

// Getenv returns the value of the specified environment variable.
// Returns the value set in `env.Env`.
func (env DummyEnv) Getenv(key string) string {
	val, ok := env.Env[key]
	if ok {
		return val
	}
	return ""
}

// HasSomeEnv returns true if one or more of the given environment variables are set.
func (env DummyEnv) HasSomeEnv(keys ...string) bool {
	for _, key := range keys {
		if env.Getenv(key) != "" {
			return true
		}
	}
	return false
}
