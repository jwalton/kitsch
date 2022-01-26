package modules

import (
	"io/fs"
	"os"
	"sync"
	"testing/fstest"

	"github.com/jwalton/kitsch/internal/cache"
	"github.com/jwalton/kitsch/internal/fileutils"
	"github.com/jwalton/kitsch/internal/gitutils"
	"github.com/jwalton/kitsch/internal/kitsch/env"
	"github.com/jwalton/kitsch/internal/kitsch/getters"
	"github.com/jwalton/kitsch/internal/kitsch/projects"
	"github.com/jwalton/kitsch/internal/kitsch/styling"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// Globals is a collection of "global" values that are passed to all modules.
// These values are available to templates via the ".Globals" property.
type Globals struct {
	// CWD is the current working directory.
	CWD string `yaml:"cwd"`
	// logicalCWD is the current working directory to display.
	logicalCWD string `yaml:"logicalCwd"`
	// Home is the user's home directory.
	Home string `yaml:"home"`
	// IsRoot is true if this is a non-windows system, and the user is UID 0.
	IsRoot bool `yaml:"isRoot"`
	// Hostname is the name of the current machine.
	Hostname string `yaml:"hostname"`
	// Jobs is the number of jobs that the shell is currently running.
	Jobs int `yaml:"jobs"`
	// Status is the return status of the previous command.
	Status int `yaml:"previousCommandStatus"`
	// PreviousCommandDuration is the duration of the previous command, in milliseconds.
	PreviousCommandDuration int64 `yaml:"previousCommandDuration"`
	// Keymap is the zsh/fish keymap. This will be "" if vi mode is not enabled,
	// "" or "main" in insert mode, and "vicmd" in normal mode.
	Keymap string `yaml:"keymap"`
	// Shell is the type of the shell (e.g. "zsh", "bash", "powershell", etc...).
	Shell string `yaml:"shell"`
	// TerminalWidth is the width of the terminal, in characters.
	TerminalWidth int `yaml:"width"`
	// PathSeparator is the path separator for the current system.
	PathSeparator string `yaml:"pathSeparator"`
}

// NewGlobals creates a new Globals object.
func NewGlobals(
	shell string,
	cwd string,
	logicalCWD string,
	terminalWidth int,
	status int,
	jobs int,
	previousCommandDuration int64,
	keymap string,
) Globals {
	var err error

	if cwd == "" {
		cwd, err = os.Getwd()
		if err != nil {
			cwd = "."
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	if terminalWidth <= 0 {
		terminalWidth, _, err = term.GetSize(0)
		if err != nil {
			terminalWidth = 80
		}
	}

	return Globals{
		CWD:                     cwd,
		logicalCWD:              logicalCWD,
		Home:                    home,
		IsRoot:                  os.Geteuid() == 0,
		Hostname:                hostname,
		Status:                  status,
		Jobs:                    jobs,
		PreviousCommandDuration: previousCommandDuration,
		Keymap:                  keymap,
		Shell:                   shell,
		TerminalWidth:           terminalWidth,
		PathSeparator:           string(os.PathSeparator),
	}
}

// LogicalCWD returns the CWD to display in the directory module.
func (globals Globals) LogicalCWD() string {
	if globals.logicalCWD == "" {
		return globals.CWD
	}
	return globals.logicalCWD
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

// GetWorkingDirectory returns the current working directory.
func (context *Context) GetWorkingDirectory() fileutils.Directory {
	return context.Directory
}

// Getenv returns the value of the specified environment variable.
func (context *Context) Getenv(key string) string {
	return context.Environment.Getenv(key)
}

// GetValueCache returns the value cache.
func (context *Context) GetValueCache() cache.Cache {
	return context.ValueCache
}

// Make sure that Context implements the GetterContext interface.
var _ getters.GetterContext = (*Context)(nil)

// Git returns a git instance for the current repo, or nil if the current
// working directory is not part of a git repo, or git is not installed.
func (context *Context) Git() gitutils.Git {
	context.mutex.Lock()
	defer context.mutex.Unlock()

	if !context.gitInitialized {
		context.git = gitutils.NewCaching("git", context.Globals.CWD)
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

// DemoConfig is a structure used to create a "demo context".  This is used for
// unit testing, and for running kitsch prompt in "demo mode" where kitsch prompt
// will not attempt to access the filesystem or environment.
type DemoConfig struct {
	// Globals are global values that are passed to all modules.
	Globals Globals `yaml:"globals"`
	// Env are environment variables to use.
	Env map[string]string `yaml:"env"`
	// Git is the git instance to use.
	Git gitutils.DemoGit `yaml:"git"`
	// CWDIsReadOnly is true if the current working directory is read-only.
	CWDIsReadOnly bool `yaml:"cwdIsReadOnly"`
}

// Load will load the demo configuration from the specified file.
func (demoConfig *DemoConfig) Load(filename string) error {
	// Set sensible defaults.
	demoConfig.Globals = Globals{
		CWD:           "/users/jwalton",
		Home:          "/users/jwalton",
		IsRoot:        false,
		Hostname:      "orac",
		Shell:         "demo",
		TerminalWidth: 80,
		PathSeparator: "/",
	}
	demoConfig.Env = map[string]string{
		"USER": "jwalton",
	}

	yamlData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// FIXME: Should do strict validation of the yaml content.

	err = yaml.Unmarshal(yamlData, demoConfig)

	return err
}

// NewDemoContext creates a demo context object.
func NewDemoContext(
	config DemoConfig,
	styles styling.Registry,
) Context {
	var cwdMode fs.FileMode = 0755
	if config.CWDIsReadOnly {
		cwdMode = 0555
	}

	demoFsys := fstest.MapFS{
		".": {Mode: cwdMode},
	}

	return Context{
		Globals:        config.Globals,
		Directory:      fileutils.NewDirectoryTestFS(config.Globals.CWD, demoFsys),
		Environment:    env.DummyEnv{Env: config.Env},
		ProjectTypes:   []projects.ProjectType{},
		ValueCache:     cache.NewMemoryCache(),
		Styles:         styles,
		gitInitialized: true,
		git:            config.Git,
	}
}

// newTestContext creates a Context with reasonable defaults that can
// be passed in to modules when unit testing.
func newTestContext(username string) *Context {
	fsys := fstest.MapFS{
		".": {Mode: 0755},
	}

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
