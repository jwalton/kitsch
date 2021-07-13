package fileutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FindFileInAncestors searches for a file or directory with the given name
// in the specified folder, or in any ancestor of that folder.  Returns the
// path to the file, or empty string if the file could not be found.
//
// This is used, e.g. to find a .git folder to check and see if we are inside
// a git repo, or to find a `package.json` file to check and see if we are
// currently in a node.js project.
func FindFileInAncestors(folder string, file string) string {
	var result string

	if filepath.Base(folder) == file {
		return folder
	}

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

// FileExists returns true if the file specified by "path" exists, and the
// parent folder can be read, false otherwise.
func FileExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// Couldn't stat the file.  :(
		return false
	}
}

// ReadFile reads the contents of a file into a string.  Returns empty string
// if the file does not exist or cannot be read.
func ReadFile(path string) string {
	content, _ := ioutil.ReadFile(path)
	text := string(content)
	return strings.TrimSpace(text)
}
