package modules

import (
	"os"
	"sync"
	"testing/fstest"

	"github.com/jwalton/kitsch-prompt/internal/cache"
	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/jwalton/kitsch-prompt/internal/gitutils"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/env"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/projects"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/styling"
)

// Globals is a collection of "global" values that are passed to all modules.
// These values are available to templates via the ".Globals" property.
type Globals struct {
	// CWD is the current wordking directory.
	CWD string
	// Home is the user's home directory.
	Home string
	// IsRoot is true if this is a non-windows system, and the user is UID 0.
	IsRoot bool
	// Hostname is the name of the current machine.
	Hostname string
	// Status is the return status of the previous command.
	Status int
	// Jobs is the number of jobs that the shell is currently running.
	Jobs int
	// PreviousCommandDuration is the duration of the previous command, in milliseconds.
	PreviousCommandDuration int64
	// Keymap is the zsh/fish keymap. This will be "" if vi mode is not enabled,
	// "" or "main" in insert mode, and "vicmd" in normal mode.
	Keymap string
	// Shell is the type of the shell (e.g. "zsh", "bash", "powershell", etc...).
	Shell string
}

// NewGlobals creates a new Globals object.
func NewGlobals(
	shell string,
	status int,
	jobs int,
	previousCommandDuration int64,
	keymap string,
) Globals {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	return Globals{
		CWD:                     cwd,
		Home:                    home,
		IsRoot:                  os.Geteuid() == 0,
		Hostname:                hostname,
		Status:                  status,
		Jobs:                    jobs,
		PreviousCommandDuration: previousCommandDuration,
		Keymap:                  keymap,
		Shell:                   shell,
	}
}

// Context is a set of common parameters passed to Module.Execute.
type Context struct {
	// Globals is a collection of "global" values that are passed to all modules.
	// These values are available to templates via the ".Globals" property.
	Globals Globals
	// Directory is the current working directory.
	Directory fileutils.Directory
	// Environment is the environment to fetch data from.
	Environment env.Env
	// ProjectTypes is the list of available project types.
	ProjectTypes []projects.ProjectType
	// The cache to retrieve values from.
	ValueCache cache.Cache
	// Styles is the style registry to use to create styles.
	Styles styling.Registry

	mutex          sync.Mutex
	gitInitialized bool
	git            gitutils.Git
}

// Git returns a git instance for the current repo, or nil if the current
// working directory is not part of a git repo, or git is not installed.
func (context *Context) Git() gitutils.Git {
	context.mutex.Lock()
	defer context.mutex.Unlock()

	if !context.gitInitialized {
		context.git = gitutils.New("git", context.Globals.CWD)
		context.gitInitialized = true
	}
	return context.git
}

// NewContext creates a new Context object for executing modules.
func NewContext(
	globals Globals,
	projectTypes []projects.ProjectType,
	cacheDir string,
	styles styling.Registry,
) Context {
	return Context{
		Globals:      globals,
		Directory:    fileutils.NewDirectory(globals.CWD),
		Environment:  env.New(),
		ProjectTypes: projectTypes,
		ValueCache:   cache.NewFileCache(cacheDir),
		Styles:       styles,
	}
}

// newTestContext creates a Context with reasonable defaults that can
// be passed in to modules when unit testing.
func newTestContext(username string) *Context {
	fsys := fstest.MapFS{}

	return &Context{
		Globals: Globals{
			CWD:                     "/Users/" + username,
			Home:                    "/Users/" + username,
			IsRoot:                  false,
			Hostname:                "lucid",
			Status:                  0,
			Jobs:                    0,
			PreviousCommandDuration: 0,
			Keymap:                  "",
			Shell:                   "bash",
		},
		Directory: fileutils.NewDirectoryTestFS("/Users/"+username, fsys),
		Environment: &env.DummyEnv{
			Env: map[string]string{
				"USER": username,
				"HOME": "/Users/" + username,
			},
		},
		ProjectTypes:   projects.DefaultProjectTypes,
		ValueCache:     cache.NewMemoryCache(),
		Styles:         styling.Registry{},
		gitInitialized: true,
		git:            nil,
	}
}
