package fileutils

import (
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// FileUtils represents utilities for interacting with the file system.
type FileUtils interface {
	// FindFileInAncestors searches for a file or directory with the given name
	// in the specified folder, or in any ancestor of that folder.  Returns the
	// path to the file, or empty string if the file could not be found.
	FindFileInAncestors(folder string, file string) string
	// FileExists returns true if the file exists on the file system.
	FileExists(path string) bool
	// ReadFile returns the contents of a file, or an error if the file cannot
	// be read.
	ReadFile(path string) ([]byte, error)
	// LookPathSafe is like exec.LookPath, but does not search in ".", event
	// if it is in the path.
	LookPathSafe(file string) (string, error)
}

type fileUtils struct{}

// New creates a new instance of FileUtils.
func New() FileUtils {
	return &fileUtils{}
}

// FindFileInAncestors searches for a file or directory with the given name
// in the specified folder, or in any ancestor of that folder.  Returns the
// path to the file, or empty string if the file could not be found.
//
// This is used, e.g. to find a .git folder to check and see if we are inside
// a git repo, or to find a `package.json` file to check and see if we are
// currently in a node.js project.
func (files *fileUtils) FindFileInAncestors(folder string, file string) string {
	var result string

	basePath := filepath.Clean(folder)

	done := false
	for !done {
		testPath := filepath.Join(basePath, file)
		if files.FileExists(testPath) {
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
func (*fileUtils) FileExists(path string) bool {
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

// ReadFile reads the contents of a file into a string.  Returns empty string
// if the file does not exist or cannot be read.
func (*fileUtils) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

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

func (*fileUtils) LookPathSafe(file string) (string, error) {
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
