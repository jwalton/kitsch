package fileutils

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// isExecutable returns an error if file is not an executable.
func isExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return fs.ErrPermission
}

// LookPathSafe is like exec.LookPath, but does not search in ".", even
// if it is in the path.
func LookPathSafe(file string) (string, error) {
	// If the file is absolute, don't try to search in the path.
	if filepath.IsAbs(file) {
		err := isExecutable(file)
		if err != nil {
			return "", &exec.Error{Name: file, Err: err}
		}
		return file, nil
	}

	// If the file contains a "/", don't try to look it up - we do not allow relative paths.
	if strings.ContainsRune(file, filepath.Separator) {
		return "", &exec.Error{Name: file, Err: fmt.Errorf("invalid path")}
	}

	// Search the path.
	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		if dir == "" || dir == "." {
			continue
		}
		path := filepath.Join(dir, file)
		if err := isExecutable(path); err == nil {
			return path, nil
		}
	}
	return "", &exec.Error{Name: file, Err: exec.ErrNotFound}
}

// FileExists returns true if the file specified by "path" exists, and the
// parent folder can be read, false otherwise.
func FileExists(path string) bool {
	// This saves us two allocs over using `os.Stat()`.`
	var fs syscall.Stat_t
	err := syscall.Stat(path, &fs)

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
