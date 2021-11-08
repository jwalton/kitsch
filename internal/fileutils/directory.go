package fileutils

import (
	"io/fs"
	"os"
	"strings"
	"sync"
)

// Directory represents a directory on the file system.  Directory has functions
// for efficiently checking if files or extensions exist in the directory.
type Directory interface {
	// Path returns the path to this directory on disk.
	Path() string
	// HasExtension returns true if the directory contains a file with the specified
	// extension (e.g. `HasExtension("gif")`).`  This function is thread safe.
	HasExtension(extension string) bool
	// HasFile returns true if the directory contains a file with the specified name.
	HasFile(name string) bool
	// HasDirectory returns true if the directory contains a directory with the specified name.
	HasDirectory(name string) bool
	// HasGlob returns true if the directory contains files which match the
	// specified glob pattern.  The pattern is the same as for `match.Match`.
	// The pattern may describe hierarchical paths like "*/*.js".
	HasGlob(glob string) bool
	// FileSystem returns an fs.FS rooted in the directory.
	FileSystem() fs.FS
	// FindFileInAncestors searches for the specified file in this directory or any
	// ancestor of this directory in the file system.  If the file is found,
	// the complete path is returned.
	FindFileInAncestors(name string) string
}

// NewDirectory creates a new Directory object for the directory at the given path.
func NewDirectory(path string) Directory {
	return &directory{
		path:       path,
		fileSystem: os.DirFS(path),
	}
}

// NewDirectoryTestFS returns a Directory instance that reads from the given FS.
// This is intended for testing.  For the returned Directory instance,
// directory.FindFileInAncestors() will never return true.
func NewDirectoryTestFS(path string, fs fs.FS) Directory {
	return &directory{
		path:         path,
		fileSystem:   fs,
		testInstance: true,
	}
}

type directory struct {
	path       string
	fileSystem fs.FS
	mutex      sync.Mutex
	// files is the set of all files in the current directory.
	files []fs.DirEntry
	// extensions is a map of all file extensions in the current directory, without
	// a leading ".".  For example, if there's a "foo.gif" in the current folder,
	// then `extensions["gif"]` will be set.
	extensions   map[string]interface{}
	testInstance bool
}

// Note that caller must have mutex.
func (dir *directory) lazyInitFiles() error {
	if dir.files == nil {
		files, err := fs.ReadDir(dir.fileSystem, ".")
		if err != nil {
			return err
		}
		dir.files = files
	}

	return nil
}

func (dir *directory) lazyInitExtensions() {
	if dir.extensions == nil {
		dir.extensions = make(map[string]interface{})

		for _, file := range dir.files {
			parts := strings.Split(file.Name(), ".")

			if len(parts[0]) == 0 {
				// Skip dotfiles.
				continue
			} else if len(parts) == 2 {
				// If there are two parts, use the second part as the extension.,
				ext := parts[1]
				if len(ext) > 0 {
					dir.extensions[ext] = nil
				}
			} else if len(parts) > 2 {
				// If there are multiple parts, like "foo.rc.js", then we want to
				// add "js" and "rc.js" to the list of extensions.
				for i := 1; i < len(parts); i++ {
					ext := strings.Join(parts[i:], ".")
					if len(ext) > 0 {
						dir.extensions[ext] = nil
					}
				}
			}
		}
	}
}

func (dir *directory) HasExtension(extension string) bool {
	// Strip the starting ".", if there is one.
	if len(extension) > 0 && extension[0] == '.' {
		extension = extension[1:]
	}

	if len(extension) == 0 {
		return false
	}

	dir.mutex.Lock()
	defer dir.mutex.Unlock()

	err := dir.lazyInitFiles()
	if err != nil {
		return false
	}

	dir.lazyInitExtensions()

	_, ok := dir.extensions[extension]
	return ok
}

func (dir *directory) Path() string {
	return dir.path
}

func (dir *directory) HasFile(name string) bool {
	_, err := fs.Stat(dir.fileSystem, name)
	return err == nil
}

func (dir *directory) HasDirectory(name string) bool {
	fileinfo, err := fs.Stat(dir.fileSystem, name)
	return err == nil && fileinfo.IsDir()
}

func (dir *directory) HasGlob(name string) bool {
	// TODO: Would be nice if this supported "**" style globs.
	files, err := fs.Glob(dir.fileSystem, name)
	if err != nil {
		return false
	}
	return len(files) > 0
}

func (dir *directory) FileSystem() fs.FS {
	return dir.fileSystem
}

func (dir *directory) FindFileInAncestors(name string) string {
	if dir.testInstance {
		return ""
	}
	return FindFileInAncestors(dir.Path(), name)
}
