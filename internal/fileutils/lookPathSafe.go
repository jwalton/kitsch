//go:build !windows
// +build !windows

package fileutils

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
