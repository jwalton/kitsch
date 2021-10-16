package modules

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/gitutils"
	"github.com/stretchr/testify/assert"
)

func TestDirectory(t *testing.T) {
	mod := DirectoryModule{}

	context := testContext("jwalton")
	context.Globals.CWD = "/tmp/test"

	result := mod.Execute(context)
	assert.Equal(t, ModuleResult{
		Data: map[string]interface{}{
			"Path":    "/tmp/test",
			"Default": "/tmp/test",
		},
		Text: "/tmp/test",
	}, result)
}

func TestHomeDirectory(t *testing.T) {
	mod := DirectoryModule{}

	context := testContext("jwalton")
	context.Globals.CWD = context.Globals.Home

	result := mod.Execute(context)
	assert.Equal(t, ModuleResult{
		Data: map[string]interface{}{
			"Path":    "~",
			"Default": "~",
		},
		Text: "~",
	}, result)
}

func TestHomeDirectorySubdirectory(t *testing.T) {
	mod := DirectoryModule{}

	context := testContext("jwalton")
	context.Globals.CWD = "/Users/jwalton/foo"

	result := mod.Execute(context)
	assert.Equal(t, ModuleResult{
		Data: map[string]interface{}{
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

	context := testContext("jwalton")
	context.Globals.CWD = "/Users/jwalton/foo/bar/baz/qux"

	assert.Equal(t, "…/bar/baz/qux", mod.Execute(context).Text)

	// Here we would normally truncate, but there's no point truncating when
	// we're just going to replace the home symbol with the truncation symbol.
	context.Globals.CWD = "/Users/jwalton/foo/bar/baz"
	assert.Equal(t, "~/foo/bar/baz", mod.Execute(context).Text)
}

func TestDirectoryTruncateToHome(t *testing.T) {
	mod := DirectoryModule{
		TruncationLength: 3,
	}

	context := testContext("jwalton")
	context.Globals.CWD = "/tmp/foo/bar/baz/qux"

	assert.Equal(t, "…/bar/baz/qux", mod.Execute(context).Text)
}

func TestDirectoryTruncateToGitRepo(t *testing.T) {
	mod := DirectoryModule{
		TruncateToRepo:   true,
		TruncationLength: 3,
	}

	context := testContext("jwalton")
	context.Globals.CWD = "/Users/jwalton/dev/kitsch-prompt/src"
	context.Environment = &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
			"HOME": "/Users/jwalton",
		},
		TestGit: &gitutils.GitUtils{RepoRoot: "/Users/jwalton/dev/kitsch-prompt"},
	}

	assert.Equal(t, "…/kitsch-prompt/src", mod.Execute(context).Text)
}
