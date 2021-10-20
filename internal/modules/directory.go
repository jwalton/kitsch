package modules

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultHomeSymbol = "~"
const defaultTruncationLength = 3
const defaultTruncationSymbol = "…"

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
	// git repo or not.
	TruncateToRepo bool `yaml:"truncateToRepo"`
	// TruncationLength is the number of parent folders that the current directory
	// should be truncated to.
	TruncationLength int `yaml:"truncationLength"`
	// TruncationSymbol will be added to the start of the string in place of any
	// paths that were removed.  Defaults to "…".
	TruncationSymbol string `yaml:"truncationSymbol"`
}

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
	isTruncated := false

	path := context.Globals.CWD

	git := context.Environment.Git()
	if mod.TruncateToRepo && git != nil && strings.HasPrefix(path, git.RepoRoot) {
		// Truncate to root of git repo if we're in a git repo.
		truncateToParts := strings.Split(git.RepoRoot, string(os.PathSeparator))
		truncateToPath := strings.Join(truncateToParts[:len(truncateToParts)-1], string(os.PathSeparator))
		path = truncationSymbol + mod.truncateToFolder(path, truncateToPath)
		isTruncated = true
	} else {
		// Truncate to the user's home directory, if we're in their home directory.
		home := context.Globals.Home
		isHome := strings.HasPrefix(path, home)

		if isHome {
			// Truncate to the home directory.
			path = defaultString(mod.HomeSymbol, defaultHomeSymbol) + mod.truncateToFolder(path, home)
			isTruncated = true
		}
	}

	// Truncate path `truncationLength`.
	if mod.TruncationLength > 0 {
		parts := strings.Split(path, string(os.PathSeparator))

		// Add one to truncationLength if isTruncated, because there's no sense truncating "~" to "…".
		truncationLength := mod.TruncationLength
		if isTruncated {
			truncationLength++
		}

		if len(parts) > truncationLength {
			parts = parts[len(parts)-mod.TruncationLength:]
			path = truncationSymbol + string(os.PathSeparator) + strings.Join(parts, string(os.PathSeparator))
		}
	}

	// TODO: Add read-only icon if read-only directory.

	data := map[string]interface{}{
		"Path": path,
	}

	return executeModule(context, mod.CommonConfig, data, mod.Style, path)
}

func init() {
	registerFactory("directory", func(node *yaml.Node) (Module, error) {
		var module DirectoryModule = DirectoryModule{
			TruncationSymbol: defaultTruncationSymbol,
			TruncationLength: defaultTruncationLength,
			TruncateToRepo:   true,
		}
		err := node.Decode(&module)
		return &module, err
	})
}
