package modules

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch-prompt/internal/gitutils"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/env"
	"github.com/stretchr/testify/assert"
)

func TestDirectory(t *testing.T) {
	mod := moduleFromYAMLMust("{type: directory}")

	context := newTestContext("jwalton")
	context.Globals.CWD = "/tmp/test"

	result := mod.Execute(context)
	assert.Equal(t, ModuleResult{
		Data: map[string]interface{}{
			"Path": "/tmp/test",
		},
		Text: "/tmp/test",
	}, result)
}

func TestHomeDirectory(t *testing.T) {
	mod := moduleFromYAMLMust("{type: directory}")

	context := newTestContext("jwalton")
	context.Globals.CWD = context.Globals.Home

	result := mod.Execute(context)
	assert.Equal(t, ModuleResult{
		Data: map[string]interface{}{
			"Path": "~",
		},
		Text: "~",
	}, result)
}

func TestHomeDirectorySubdirectory(t *testing.T) {
	mod := moduleFromYAMLMust("{type: directory}")

	context := newTestContext("jwalton")
	context.Globals.CWD = "/Users/jwalton/foo"

	result := mod.Execute(context)
	assert.Equal(t, ModuleResult{
		Data: map[string]interface{}{
			"Path": "~/foo",
		},
		Text: "~/foo",
	}, result)
}

func TestHomeDirectorySubdirectoryTruncated(t *testing.T) {
	mod := moduleFromYAMLMust(heredoc.Doc(`
		type: directory
		truncationLength: 3
	`))

	context := newTestContext("jwalton")
	context.Globals.CWD = "/Users/jwalton/foo/bar/baz/qux"

	assert.Equal(t, "…/bar/baz/qux", mod.Execute(context).Text)

	// Here we would normally truncate, but there's no point truncating when
	// we're just going to replace the home symbol with the truncation symbol.
	context.Globals.CWD = "/Users/jwalton/foo/bar/baz"
	assert.Equal(t, "~/foo/bar/baz", mod.Execute(context).Text)
}

func TestDirectoryTruncateToHome(t *testing.T) {
	mod := moduleFromYAMLMust(heredoc.Doc(`
		type: directory
		truncationLength: 3
	`))

	context := newTestContext("jwalton")
	context.Globals.CWD = "/tmp/foo/bar/baz/qux"

	assert.Equal(t, "…/bar/baz/qux", mod.Execute(context).Text)
}

func TestDirectoryTruncateToGitRepo(t *testing.T) {
	mod := moduleFromYAMLMust(heredoc.Doc(`
		type: directory
		truncationLength: 3
	`))

	context := newTestContext("jwalton")
	context.Environment = &env.DummyEnv{
		Env: map[string]string{
			"USER": "jwalton",
			"HOME": "/Users/jwalton",
		},
	}

	// FIXME: Create a test git instance.  This won't work unless git is installed
	// and this is in the right directory.
	context.gitInitialized = true
	context.git = gitutils.DemoGit{
		RepoRootDirectory:     "/Users/jwalton/dev/kitsch-prompt",
		Head:                  "master",
		IsDetached:            false,
		CurrentBranchUpstream: "master",
		CurrentState:          gitutils.StateNone,
	}

	context.Globals.CWD = "/Users/jwalton/dev/kitsch-prompt/src"
	assert.Equal(t, "…/kitsch-prompt/src", mod.Execute(context).Text)

	context.Globals.CWD = "/Users/jwalton/dev/kitsch-prompt/src/foo/bar/baz/qux"
	assert.Equal(t, "…/bar/baz/qux", mod.Execute(context).Text)

	mod = moduleFromYAMLMust(heredoc.Doc(`
		type: directory
		truncationLength: 3
		truncateToRepo: false
	`))

	context.Globals.CWD = "/Users/jwalton/dev/kitsch-prompt/src"
	assert.Equal(t, "~/dev/kitsch-prompt/src", mod.Execute(context).Text)

	context.Globals.CWD = "/Users/jwalton/work/dev/kitsch-prompt/src"
	assert.Equal(t, "…/dev/kitsch-prompt/src", mod.Execute(context).Text)

}
