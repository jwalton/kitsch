package modules

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultHomeSymbol = "~"
const defaultTruncationLength = 3
const defaultTruncationSymbol = "…"

func getVolumeName(path string) string {
	return filepath.VolumeName(path)
}

// DirectoryModule shows the current working directory.
//
// Provides the following template variables:
//
// • Path - The directory to show.
//
type DirectoryModule struct {
	CommonConfig `yaml:",inline"`
	// HomeSymbol is the symbol to replace the home directory with in directory
	// strings.  Defaults to "~".
	HomeSymbol string `yaml:"homeSymbol"`
	// TruncateToRepo controls whether we truncate to the root directory of the
	// git repo or not.  If this is true, and we are in a source code repository,
	// we will replace everything up to the repo root directory with RepoSymbol.
	TruncateToRepo bool `yaml:"truncateToRepo"`
	// RepoSymbol is a string that will be added as a prefix when we truncate to a repo.
	RepoSymbol string `yaml:"repoSymbol"`
	// TruncationLength is the number of parent folders that the current directory
	// should be truncated to.  If 0, truncation will be disabled.
	TruncationLength int `yaml:"truncationLength"`
	// TruncationSymbol will be added to the start of the string in place of any
	// paths that were removed.  Defaults to "…".
	TruncationSymbol string `yaml:"truncationSymbol"`

	getVolumeName func(string) string
}

type directoryModuleResult struct {
	// Path is the path that will be shown to the user.
	Path string
	// PathSeparator is the system defined path separator.
	PathSeparator string
}

// Removes `truncatePath` from the start of `path`.  The returned path will always
// start with a path separator.
func (mod DirectoryModule) truncateToFolder(path string, truncatePath string) string {
	charsToStrip := len(truncatePath)

	truncatePathEndsWithSeparator := len(truncatePath) > 0 && truncatePath[len(truncatePath)-1] == os.PathSeparator
	if truncatePathEndsWithSeparator {
		// Leave the separator
		charsToStrip--
	}

	return path[charsToStrip:]
}

// Execute the directory module.
func (mod DirectoryModule) Execute(context *Context) ModuleResult {
	truncationSymbol := defaultString(mod.TruncationSymbol, defaultTruncationSymbol)

	pathSeparator := context.Globals.PathSeparator
	path := context.Globals.CWD
	volumeName := mod.getVolumeName(path)
	isHome := strings.HasPrefix(path, context.Globals.Home)

	prefix := ""
	git := context.Git()
	if mod.TruncateToRepo && git != nil && strings.HasPrefix(path, git.RepoRoot()) {
		// Truncate to root of git repo if we're in a git repo.
		gitRepoParts := strings.Split(git.RepoRoot(), pathSeparator)
		prefix = mod.RepoSymbol + gitRepoParts[len(gitRepoParts)-1]
		path = mod.truncateToFolder(path, git.RepoRoot())
		isHome = false
	} else if volumeName != "" && !isHome {
		// If the path starts with a volume name, remove it.
		prefix = volumeName
		path = path[len(volumeName):]
	}

	// Truncate path `truncationLength`.
	if mod.TruncationLength > 0 {
		truncationLength := mod.TruncationLength
		if isHome {
			// Truncate to the user's home directory, if we're in their home directory.
			path = mod.HomeSymbol + mod.truncateToFolder(path, context.Globals.Home)

			// Add one to the truncation length, as there's no sense in replaceing
			// "~" with "..."
			truncationLength++
		}

		parts := strings.Split(path, pathSeparator)

		if len(parts) > truncationLength {
			parts = parts[len(parts)-mod.TruncationLength:]
			path = truncationSymbol + pathSeparator + strings.Join(parts, pathSeparator)
			// If there's a prefix, add a leading path separator before the truncation symbol.
			if prefix != "" {
				path = pathSeparator + path
			}
		}
	}

	// TODO: Add read-only icon if read-only directory.

	data := directoryModuleResult{
		Path:          prefix + path,
		PathSeparator: pathSeparator,
	}

	return executeModule(context, mod.CommonConfig, data, mod.Style, data.Path)
}

func init() {
	registerFactory("directory", func(node *yaml.Node) (Module, error) {
		var module DirectoryModule = DirectoryModule{
			HomeSymbol:       defaultHomeSymbol,
			TruncationSymbol: defaultTruncationSymbol,
			TruncationLength: defaultTruncationLength,
			TruncateToRepo:   true,
			getVolumeName:    getVolumeName,
		}
		err := node.Decode(&module)
		return &module, err
	})
}
