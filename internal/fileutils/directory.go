package fileutils

import (
	"io"
	"io/fs"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jwalton/kitsch/internal/kitsch/log"
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
	// Stat returns the os.FileInfo for the specified file.
	Stat(path string) (os.FileInfo, error)
}

// NewDirectory creates a new Directory object for the directory at the given path.
// scanTimeout is the maximum time to spend reading files from disk.  If 0,
// no timeout is used.
func NewDirectory(path string, scanTimeout time.Duration) Directory {
	return &directory{
		path:         path,
		fileSystem:   os.DirFS(path),
		testInstance: false,
		scanTimeout:  scanTimeout,
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
		scanTimeout:  0,
	}
}

type directory struct {
	path       string
	fileSystem fs.FS
	filesOnce  sync.Once
	// files is a map of all files in the current directory.
	files map[string]interface{}
	// extensions is a map of all file extensions in the current directory, without
	// a leading ".".  For example, if there's a "foo.gif" in the current folder,
	// then `extensions["gif"]` will be set.
	extensions   map[string]interface{}
	testInstance bool
	// scanTimeout is the maximum time to wait for loading directory contents to complete.
	scanTimeout time.Duration
}

// Note that caller must have mutex.
func (dir *directory) lazyInitFiles() {
	dir.filesOnce.Do(func() {
		dir.files = make(map[string]interface{})
		dir.extensions = make(map[string]interface{})

		if dir.testInstance {
			// For test instance, use fs.ReadDir.  This is MUCH
			// slower, so we don't use it in "production", but
			// it does let us use an in-memory filesystem for tests.
			files, err := fs.ReadDir(dir.fileSystem, ".")
			if err != nil {
				return
			}

			for _, file := range files {
				dir.cacheFile(file.Name())
			}

		} else {
			f, err := os.Open(dir.path)
			if err != nil {
				return
			}
			defer f.Close()

			start := time.Now()

			done := false
			fileCount := 0
			for !done {
				// Grab files, 256 at a time.
				files, err := f.Readdirnames(256)
				if err == io.EOF {
					done = true
				} else if err != nil {
					return
				} else {
					for _, file := range files {
						fileCount++
						dir.cacheFile(file)
					}
				}

				if dir.scanTimeout > 0 && time.Since(start) > dir.scanTimeout {
					// We've spent too long reading files from the disk.  Use what
					// we have so far.
					done = true
					log.Info("Directory scan timed out ", dir.path, " after ", fileCount, " files.")
				}
			}
		}
	})
}

// cacheFile caches information about a file in the current directory.
func (dir *directory) cacheFile(filename string) {
	dir.files[filename] = nil

	parts := strings.Split(filename, ".")

	if len(parts[0]) == 0 {
		// Skip dotfiles.
	} else if len(parts) == 2 {
		// If there are two parts, use the second part as the extension.
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

func (dir *directory) HasExtension(extension string) bool {
	// Strip the starting ".", if there is one.
	if len(extension) > 0 && extension[0] == '.' {
		extension = extension[1:]
	}

	if len(extension) == 0 {
		return false
	}

	dir.lazyInitFiles()

	_, ok := dir.extensions[extension]
	return ok
}

func (dir *directory) Path() string {
	return dir.path
}

func (dir *directory) HasFile(name string) bool {
	dir.lazyInitFiles()

	_, ok := dir.files[name]
	return ok
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

func (dir *directory) Stat(path string) (os.FileInfo, error) {
	return fs.Stat(dir.fileSystem, path)
}
