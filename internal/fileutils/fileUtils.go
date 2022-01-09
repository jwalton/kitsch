package fileutils

import (
	"os"
	"path/filepath"
)

// FileExists returns true if the file specified by "path" exists, and the
// parent folder can be read, false otherwise.
func FileExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// Either `os.IsNotExist(err)`` is true, and the file doesn't exist,
		// or we couldn't state the file.  Either way, just claim it isn't there.
		return false
	}
}

// FindFileInAncestors searches for a file or directory with the given name
// in the specified folder, or in any ancestor of that folder.  Returns the
// path to the file, or empty string if the file could not be found.
func FindFileInAncestors(folder string, file string) string {
	var result string

	basePath := filepath.Clean(folder)

	done := false
	for !done {
		testPath := filepath.Join(basePath, file)
		if FileExists(testPath) {
			result = testPath
			break
		}

		newBasePath := filepath.Dir(basePath)
		if newBasePath == basePath {
			done = true
		}
		basePath = newBasePath
	}

	return result
}
