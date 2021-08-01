package modules

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/stretchr/testify/assert"
)

func TestDirectory(t *testing.T) {
	mod := DirectoryModule{}

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
			"HOME": "/Users/jwalton",
		},
		CWD: "/tmp/test",
	}

	result := mod.Execute(env)
	assert.Equal(t, ModuleResult{
		Data: map[string]interface{}{
			"rawDirectory": "/tmp/test",
			"directory":    "/tmp/test",
		},
		Text: "/tmp/test",
	}, result)
}

func TestHomeDirectory(t *testing.T) {
	mod := DirectoryModule{}

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
			"HOME": "/Users/jwalton",
		},
		CWD: "/Users/jwalton",
	}

	result := mod.Execute(env)
	assert.Equal(t, ModuleResult{
		Data: map[string]interface{}{
			"rawDirectory": "/Users/jwalton",
			"directory":    "~",
		},
		Text: "~",
	}, result)
}

func TestHomeDirectorySubdirectory(t *testing.T) {
	mod := DirectoryModule{}

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
			"HOME": "/Users/jwalton",
		},
		CWD: "/Users/jwalton/foo",
	}

	result := mod.Execute(env)
	assert.Equal(t, ModuleResult{
		Data: map[string]interface{}{
			"rawDirectory": "/Users/jwalton/foo",
			"directory":    "~/foo",
		},
		Text: "~/foo",
	}, result)
}

func TestHomeDirectorySubdirectoryTruncated(t *testing.T) {
	mod := DirectoryModule{}

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
			"HOME": "/Users/jwalton",
		},
		CWD: "/Users/jwalton/foo/bar/baz/qux",
	}

	assert.Equal(t, "…/bar/baz/qux", mod.Execute(env).Text)

	// Here we would normally truncate, but there's no point truncating when
	// we're just going to replace the home symbol with the truncation symbol.
	env.CWD = "/Users/jwalton/foo/bar/baz"
	assert.Equal(t, "~/foo/bar/baz", mod.Execute(env).Text)
}

func TestHomeDirectoryTruncated(t *testing.T) {
	mod := DirectoryModule{}

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
			"HOME": "/Users/jwalton",
		},
		CWD: "/tmp/foo/bar/baz/qux",
	}

	assert.Equal(t, "…/bar/baz/qux", mod.Execute(env).Text)
}
