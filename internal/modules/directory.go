package modules

import (
	"os"
	"strings"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"gopkg.in/yaml.v3"
)

const defaultHomeSymbol = "~"
const defaultTruncationLength = 3
const defaultTruncationSymbol = "…"

// DirectoryModule shows the current working directory.
//
// Provides the following template variables:
//
// • rawDirectory - The unmodified current working directory.
//
// • directory - The directory to show.
//
type DirectoryModule struct {
	CommonConfig `yaml:",inline"`
	// HomeSymbol is the symbol to replace the home directory with in directory
	// strings.  Defaults to "~".
	HomeSymbol string
	// TruncateToRepo controls whether we truncate to the root directory of the
	// git repo or not.
	TruncateToRepo bool
	// TruncationLength is the number of parent folders that the current directory
	// should be truncated to.
	TruncationLength int
	// TruncationSymbol will be added to the start of the string in place of any
	// paths that were removed.  Defaults to "…".
	TruncationSymbol string
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
func (mod DirectoryModule) Execute(env env.Env) ModuleResult {
	truncationSymbol := defaultString(mod.TruncationSymbol, defaultTruncationSymbol)
	isTruncated := false

	rawDirectory := env.Getwd()

	directory := rawDirectory

	git := env.Git()
	if mod.TruncateToRepo && git != nil && strings.HasPrefix(directory, git.RepoRoot) {
		// Truncate to root of git repo if we're in a git repo.
		truncateToParts := strings.Split(git.RepoRoot, string(os.PathSeparator))
		truncateToPath := strings.Join(truncateToParts[:len(truncateToParts)-1], string(os.PathSeparator))
		directory = truncationSymbol + mod.truncateToFolder(directory, truncateToPath)
		isTruncated = true
	} else {
		// Truncate to the user's home directory, if we're in their home directory.
		home := env.UserHomeDir()
		isHome := strings.HasPrefix(directory, home)

		if isHome {
			// Truncate to the home directory.
			directory = defaultString(mod.HomeSymbol, defaultHomeSymbol) + mod.truncateToFolder(directory, home)
			isTruncated = true
		}
	}

	// Truncate path `truncationLength`.
	if mod.TruncationLength > 0 {
		parts := strings.Split(directory, string(os.PathSeparator))

		// Add one to truncationLength if isTruncated, because there's no sense truncating "~" to "…".
		truncationLength := mod.TruncationLength
		if isTruncated {
			truncationLength++
		}

		if len(parts) > truncationLength {
			parts = parts[len(parts)-mod.TruncationLength:]
			directory = truncationSymbol + string(os.PathSeparator) + strings.Join(parts, string(os.PathSeparator))
		}
	}

	// TODO: Add read-only icon if read-only directory.

	data := map[string]interface{}{
		"rawDirectory": rawDirectory,
		"directory":    directory,
	}

	return executeModule(mod.CommonConfig, data, mod.Style, directory)
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
