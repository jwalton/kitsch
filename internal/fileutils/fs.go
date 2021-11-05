package fileutils

import "io/fs"

// FSFileExists returns true if the file exists in the given `fs` instance.
func FSFileExists(files fs.FS, path string) bool {
	_, err := fs.Stat(files, path)
	return err == nil
}
