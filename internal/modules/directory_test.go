package modules

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/gitutils"
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
			"RawPath": "/tmp/test",
			"Path":    "/tmp/test",
			"Default": "/tmp/test",
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
			"RawPath": "/Users/jwalton",
			"Path":    "~",
			"Default": "~",
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
			"RawPath": "/Users/jwalton/foo",
			"Path":    "~/foo",
			"Default": "~/foo",
		},
		Text: "~/foo",
	}, result)
}

func TestHomeDirectorySubdirectoryTruncated(t *testing.T) {
	mod := DirectoryModule{
		TruncationLength: 3,
	}

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

func TestDirectoryTruncateToHome(t *testing.T) {
	mod := DirectoryModule{
		TruncationLength: 3,
	}

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
			"HOME": "/Users/jwalton",
		},
		CWD: "/tmp/foo/bar/baz/qux",
	}

	assert.Equal(t, "…/bar/baz/qux", mod.Execute(env).Text)
}

func TestDirectoryTruncateToGitRepo(t *testing.T) {
	mod := DirectoryModule{
		TruncateToRepo:   true,
		TruncationLength: 3,
	}

	env := &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
			"HOME": "/Users/jwalton",
		},
		CWD:     "/Users/jwalton/dev/kitsch-prompt/src",
		TestGit: &gitutils.GitUtils{RepoRoot: "/Users/jwalton/dev/kitsch-prompt"},
	}

	assert.Equal(t, "…/kitsch-prompt/src", mod.Execute(env).Text)
}
