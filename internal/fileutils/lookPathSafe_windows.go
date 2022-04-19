//go:build windows
// +build windows

package fileutils

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func chkStat(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if d.IsDir() {
		return fs.ErrPermission
	}
	return nil
}

func hasExt(file string) bool {
	i := strings.LastIndex(file, ".")
	if i < 0 {
		return false
	}
	return strings.LastIndexAny(file, `:\/`) < i
}

func findExecutable(file string, exts []string) (string, error) {
	if len(exts) == 0 {
		return file, chkStat(file)
	}
	if hasExt(file) {
		if chkStat(file) == nil {
			return file, nil
		}
	}
	for _, e := range exts {
		if f := file + e; chkStat(f) == nil {
			return f, nil
		}
	}
	return "", fs.ErrNotExist
}

// LookPathSafe is like exec.LookPath, but does not search in ".", even
// if it is in the path.
func LookPathSafe(file string) (string, error) {
	var exts []string
	x := os.Getenv(`PATHEXT`)
	if x != "" {
		for _, e := range strings.Split(strings.ToLower(x), `;`) {
			if e == "" {
				continue
			}
			if e[0] != '.' {
				e = "." + e
			}
			exts = append(exts, e)
		}
	} else {
		exts = []string{".com", ".exe", ".bat", ".cmd"}
	}

	if strings.ContainsAny(file, `:\/`) {
		f, err := findExecutable(file, exts)
		if err != nil {
			return "", &exec.Error{Name: file, Err: err}
		}
		return f, nil
	}
	if f, err := findExecutable(filepath.Join(".", file), exts); err == nil {
		return f, nil
	}
	path := os.Getenv("path")
	for _, dir := range filepath.SplitList(path) {
		if dir == "." {
			continue
		}
		if f, err := findExecutable(filepath.Join(dir, file), exts); err == nil {
			return f, nil
		}
	}
	return "", &exec.Error{Name: file, Err: exec.ErrNotFound}
}
