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
	// TruncationLength is the number of parent folders that the current directory
	// should be truncated to.
	TruncationLength int
	// TruncationSymbol will be added to the start of the string in place of any
	// paths that were removed.  Defaults to "…".
	TruncationSymbol string
}

// Execute the directory module.
func (mod DirectoryModule) Execute(env env.Env) ModuleResult {
	rawDirectory := env.Getwd()

	// TODO: Truncate to root of git repo if we're in a git repo.
	directory := rawDirectory

	home := env.UserHomeDir()
	isHome := strings.HasPrefix(directory, home)

	if isHome {
		// Strip the leading home directory
		homeHasSeparator := len(home) > 0 && home[len(home)-1] == os.PathSeparator
		if homeHasSeparator {
			// Make sure we leave the leading separator in "directory".
			directory = directory[len(home)-1:]
		} else {
			directory = directory[len(home):]
		}
	}

	truncationLength := defaultNumber(mod.TruncationLength, defaultTruncationLength)

	// Truncate directory to `truncationLength`.
	//
	// Note if `isHome`, we add one to truncationLength, because there's no sense
	// truncating "~" to "…".
	parts := strings.Split(directory, string(os.PathSeparator))
	isTruncated := isHome && (len(parts) > truncationLength+1) || !isHome && (len(parts) > truncationLength)
	if isTruncated {
		parts = parts[len(parts)-truncationLength:]
		truncationSymbol := defaultString(mod.TruncationSymbol, defaultTruncationSymbol)
		directory = truncationSymbol + string(os.PathSeparator) + strings.Join(parts, string(os.PathSeparator))
	}

	// If this is the home directory, and we haven't truncated the path, add the
	// home symbol back.
	if isHome && !isTruncated {
		directory = defaultString(mod.HomeSymbol, defaultHomeSymbol) + directory
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
		var module DirectoryModule
		err := node.Decode(&module)
		return &module, err
	})
}
