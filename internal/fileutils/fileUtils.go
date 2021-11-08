package fileutils

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func findExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return fs.ErrPermission
}

// LookPathSafe is like exec.LookPath, but does not search in ".", event
// if it is in the path.
func LookPathSafe(file string) (string, error) {
	if strings.Contains(file, "/") {
		err := findExecutable(file)
		if err == nil {
			return file, nil
		}
		return "", &exec.Error{Name: file, Err: err}
	}

	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		if dir == "" || dir == "." {
			continue
		}
		path := filepath.Join(dir, file)
		if err := findExecutable(path); err == nil {
			return path, nil
		}
	}
	return "", &exec.Error{Name: file, Err: exec.ErrNotFound}
}

// fileExists returns true if the file specified by "path" exists, and the
// parent folder can be read, false otherwise.
func fileExists(path string) bool {
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
		if fileExists(testPath) {
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
