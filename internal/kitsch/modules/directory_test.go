package modules

import (
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/jwalton/kitsch-prompt/internal/gitutils"
	"github.com/stretchr/testify/assert"
)

func makeTestDirectoryModule(
	pathSeparator string,
	cwd string,
	gitRoot string,
	yaml string,
) (*Context, *DirectoryModule) {
	context := newTestContext("jwalton")
	context.Globals.CWD = cwd
	context.Globals.PathSeparator = pathSeparator

	fsys := fstest.MapFS{
		".": {Mode: 0755},
	}
	fileutils.NewDirectoryTestFS(cwd, fsys)

	// Setup git
	context.gitInitialized = true
	if gitRoot != "" {
		context.git = gitutils.DemoGit{
			RepoRootDirectory:     "/Users/jwalton/dev/kitsch-prompt",
			Head:                  "master",
			IsDetached:            false,
			CurrentBranchUpstream: "master",
			CurrentState:          gitutils.StateNone,
		}
	}

	mod := moduleFromYAMLMust(yaml).(*DirectoryModule)
	mod.getVolumeName = func(path string) string {
		if pathSeparator != "/" && len(path) >= 3 && path[1] == ':' && string(path[2]) == pathSeparator {
			return path[:2]
		}
		return ""
	}

	return context, mod
}

func TestDirectory(t *testing.T) {
	context, mod := makeTestDirectoryModule("/", "/tmp/test", "", "{type: directory}")

	result := mod.Execute(context)
	assert.Equal(t, ModuleResult{
		Data: directoryModuleResult{
			Path:           "/tmp/test",
			PathSeparator:  "/",
			ReadOnly:       false,
			ReadOnlySymbol: "🔒",
		},
		Text: "/tmp/test",
	}, result)
}

func TestReadOnlyDirectory(t *testing.T) {
	context, mod := makeTestDirectoryModule("/", "/tmp/test", "", "{type: directory}")
	fsys := fstest.MapFS{
		".": {Mode: 0555},
	}
	context.Directory = fileutils.NewDirectoryTestFS(context.Globals.CWD, fsys)

	result := mod.Execute(context)
	assert.Equal(t, ModuleResult{
		Data: directoryModuleResult{
			Path:           "/tmp/test",
			PathSeparator:  "/",
			ReadOnly:       true,
			ReadOnlySymbol: "🔒",
		},
		Text: "/tmp/test🔒",
	}, result)
}

func TestRootDirectory(t *testing.T) {
	context, mod := makeTestDirectoryModule("/", "/", "", "{type: directory}")
	assert.Equal(t, "/", mod.Execute(context).Text)
}

func TestRootDirectoryWindows(t *testing.T) {
	context, mod := makeTestDirectoryModule("\\", "D:\\", "", "{type: directory}")
	assert.Equal(t, "D:\\", mod.Execute(context).Text)
}

func TestHomeDirectory(t *testing.T) {
	context, mod := makeTestDirectoryModule("/", "/Users/jwalton", "", "{type: directory}")
	context.Globals.CWD = context.Globals.Home
	assert.Equal(t, "~", mod.Execute(context).Text)

	context.Globals.CWD = context.Globals.Home + "/foo"
	assert.Equal(t, "~/foo", mod.Execute(context).Text)

	context.Globals.CWD = context.Globals.Home + "/foo/bar/baz"
	assert.Equal(t, "~/foo/bar/baz", mod.Execute(context).Text)

	context.Globals.CWD = context.Globals.Home + "/foo/bar/baz/qux"
	assert.Equal(t, "…/bar/baz/qux", mod.Execute(context).Text)
}

func TestHomeDirectoryWindows(t *testing.T) {
	context, mod := makeTestDirectoryModule("\\", "C:\\Users\\jwalton", "", "{type: directory}")
	context.Globals.Home = "C:\\Users\\jwalton"
	assert.Equal(t, "~", mod.Execute(context).Text)

	context.Globals.CWD = context.Globals.Home + "\\foo"
	assert.Equal(t, "~\\foo", mod.Execute(context).Text)

	context.Globals.CWD = context.Globals.Home + "\\foo\\bar\\baz"
	assert.Equal(t, "~\\foo\\bar\\baz", mod.Execute(context).Text)

	context.Globals.CWD = context.Globals.Home + "\\foo\\bar\\baz\\qux"
	assert.Equal(t, "…\\bar\\baz\\qux", mod.Execute(context).Text)
}

func TestDirectoryTruncate(t *testing.T) {
	context, mod := makeTestDirectoryModule("/", "/tmp/foo/bar/baz/qux", "",
		heredoc.Doc(`
			type: directory
			truncationLength: 3
		`),
	)

	assert.Equal(t, "…/bar/baz/qux", mod.Execute(context).Text)
}

func TestDirectoryTruncateWindows(t *testing.T) {
	context, mod := makeTestDirectoryModule("\\", "C:\\tmp\\foo\\bar\\baz\\qux", "",
		heredoc.Doc(`
			type: directory
			truncationLength: 3
		`),
	)

	assert.Equal(t, "C:\\…\\bar\\baz\\qux", mod.Execute(context).Text)
}

func TestDirectoryTruncateToGitRepo(t *testing.T) {
	context, mod := makeTestDirectoryModule("/", "/Users/jwalton/dev/kitsch-prompt", "/Users/jwalton/dev/kitsch-prompt",
		heredoc.Doc(`
			type: directory
			truncationLength: 3
		`),
	)
	assert.Equal(t, "kitsch-prompt", mod.Execute(context).Text)

	context.Globals.CWD = "/Users/jwalton/dev/kitsch-prompt/src"
	assert.Equal(t, "kitsch-prompt/src", mod.Execute(context).Text)

	context.Globals.CWD = "/Users/jwalton/dev/kitsch-prompt/src/foo/bar/baz/qux"
	assert.Equal(t, "kitsch-prompt/…/bar/baz/qux", mod.Execute(context).Text)

	mod.RepoSymbol = "?"
	assert.Equal(t, "?kitsch-prompt/…/bar/baz/qux", mod.Execute(context).Text)

	context, mod = makeTestDirectoryModule("/", "/Users/jwalton/dev/kitsch-prompt", "/Users/jwalton/dev/kitsch-prompt",
		heredoc.Doc(`
			type: directory
			truncationLength: 3
			truncateToRepo: false
		`),
	)
	context.Globals.CWD = "/Users/jwalton/dev/kitsch-prompt/src"
	assert.Equal(t, "~/dev/kitsch-prompt/src", mod.Execute(context).Text)

	context.Globals.CWD = "/Users/jwalton/work/dev/kitsch-prompt/src"
	assert.Equal(t, "…/dev/kitsch-prompt/src", mod.Execute(context).Text)

}
