package modules

import (
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch/internal/fileutils"
	"github.com/jwalton/kitsch/internal/gitutils"
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

	// Setup git
	context.gitInitialized = true
	if gitRoot != "" {
		context.git = gitutils.DemoGit{
			RepoRootDirectory:     "/Users/jwalton/dev/kitsch",
			HeadDescription:       "master",
			IsDetached:            false,
			CurrentBranchUpstream: "master",
			CurrentState:          gitutils.StateNone,
		}
	}

	mod := moduleFromYAML(yaml).(*DirectoryModule)
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
		DefaultText: "/tmp/test",
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
		DefaultText: "/tmp/test🔒",
	}, result)
}

func TestRootDirectory(t *testing.T) {
	context, mod := makeTestDirectoryModule("/", "/", "", "{type: directory}")
	assert.Equal(t, "/", mod.Execute(context).DefaultText)
}

func TestRootDirectoryWindows(t *testing.T) {
	context, mod := makeTestDirectoryModule("\\", "D:\\", "", "{type: directory}")
	assert.Equal(t, "D:\\", mod.Execute(context).DefaultText)
}

func TestHomeDirectory(t *testing.T) {
	context, mod := makeTestDirectoryModule("/", "/Users/jwalton", "", "{type: directory}")
	context.Globals.CWD = context.Globals.Home
	assert.Equal(t, "~", mod.Execute(context).DefaultText)

	context.Globals.CWD = context.Globals.Home + "/foo"
	assert.Equal(t, "~/foo", mod.Execute(context).DefaultText)

	context.Globals.CWD = context.Globals.Home + "/foo/bar/baz"
	assert.Equal(t, "~/foo/bar/baz", mod.Execute(context).DefaultText)

	context.Globals.CWD = context.Globals.Home + "/foo/bar/baz/qux"
	assert.Equal(t, "…/bar/baz/qux", mod.Execute(context).DefaultText)
}

func TestHomeDirectoryWindows(t *testing.T) {
	context, mod := makeTestDirectoryModule("\\", "C:\\Users\\jwalton", "", "{type: directory}")
	context.Globals.Home = "C:\\Users\\jwalton"
	assert.Equal(t, "~", mod.Execute(context).DefaultText)

	context.Globals.CWD = context.Globals.Home + "\\foo"
	assert.Equal(t, "~\\foo", mod.Execute(context).DefaultText)

	context.Globals.CWD = context.Globals.Home + "\\foo\\bar\\baz"
	assert.Equal(t, "~\\foo\\bar\\baz", mod.Execute(context).DefaultText)

	context.Globals.CWD = context.Globals.Home + "\\foo\\bar\\baz\\qux"
	assert.Equal(t, "…\\bar\\baz\\qux", mod.Execute(context).DefaultText)
}

func TestDirectoryTruncate(t *testing.T) {
	context, mod := makeTestDirectoryModule("/", "/tmp/foo/bar/baz/qux", "",
		heredoc.Doc(`
			type: directory
			truncationLength: 3
		`),
	)

	assert.Equal(t, "…/bar/baz/qux", mod.Execute(context).DefaultText)
}

func TestDirectoryTruncateWindows(t *testing.T) {
	context, mod := makeTestDirectoryModule("\\", "C:\\tmp\\foo\\bar\\baz\\qux", "",
		heredoc.Doc(`
			type: directory
			truncationLength: 3
		`),
	)

	assert.Equal(t, "C:\\…\\bar\\baz\\qux", mod.Execute(context).DefaultText)

	context.Globals.CWD = "C:\\tmp\\foo\\bar\\baz"
	assert.Equal(t, "C:\\…\\foo\\bar\\baz", mod.Execute(context).DefaultText)

	context.Globals.CWD = "C:\\tmp\\foo\\bar"
	assert.Equal(t, "C:\\tmp\\foo\\bar", mod.Execute(context).DefaultText)
}

func TestDirectoryTruncateToGitRepo(t *testing.T) {
	context, mod := makeTestDirectoryModule("/", "/Users/jwalton/dev/kitsch", "/Users/jwalton/dev/kitsch",
		heredoc.Doc(`
			type: directory
			truncationLength: 3
		`),
	)
	assert.Equal(t, "kitsch", mod.Execute(context).DefaultText)

	context.Globals.CWD = "/Users/jwalton/dev/kitsch/src"
	assert.Equal(t, "kitsch/src", mod.Execute(context).DefaultText)

	context.Globals.CWD = "/Users/jwalton/dev/kitsch/src/foo/bar/baz/qux"
	assert.Equal(t, "kitsch/…/bar/baz/qux", mod.Execute(context).DefaultText)

	mod.RepoSymbol = "?"
	assert.Equal(t, "?kitsch/…/bar/baz/qux", mod.Execute(context).DefaultText)

	context, mod = makeTestDirectoryModule("/", "/Users/jwalton/dev/kitsch", "/Users/jwalton/dev/kitsch",
		heredoc.Doc(`
			type: directory
			truncationLength: 3
			truncateToRepo: false
		`),
	)
	context.Globals.CWD = "/Users/jwalton/dev/kitsch/src"
	assert.Equal(t, "~/dev/kitsch/src", mod.Execute(context).DefaultText)

	context.Globals.CWD = "/Users/jwalton/work/dev/kitsch/src"
	assert.Equal(t, "…/dev/kitsch/src", mod.Execute(context).DefaultText)

}

func TestWindowLogicalCWD(t *testing.T) {
	context, mod := makeTestDirectoryModule("\\", "C:\\Users\\jwalton", "", "{type: directory}")
	context.Globals.Home = "C:\\Users\\jwalton"

	context.Globals.logicalCWD = "Env:\\"

	assert.Equal(t, "Env:\\", mod.Execute(context).DefaultText)
}
