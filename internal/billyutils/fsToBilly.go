package billyutils

import (
	"io"
	"io/fs"
	"path"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
)

// FsToBilly converts a fs.FS filesystem to a billy.Filesystem.
func FsToBilly(in fs.FS) (billy.Filesystem, error) {
	memory := memfs.New()
	err := copyToBilly(in, memory, ".")
	if err != nil {
		return nil, err
	}
	return memory, err
}

func copyToBilly(in fs.FS, out billy.Filesystem, dir string) error {
	billyDir := dir
	if billyDir == "." {
		billyDir = "/"
	}

	files, err := fs.ReadDir(in, dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := path.Join(dir, file.Name())
		billyFileName := path.Join(billyDir, file.Name())

		if file.IsDir() {
			info, err := file.Info()
			if err != nil {
				return err
			}

			err = out.MkdirAll(billyFileName, info.Mode())
			if err != nil {
				return err
			}
			err = copyToBilly(in, out, fileName)
			if err != nil {
				return err
			}

			continue
		}

		src, err := in.Open(fileName)
		if err != nil {
			return err
		}

		dst, err := out.Create(billyFileName)
		if err != nil {
			return err
		}

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		if err := dst.Close(); err != nil {
			return err
		}

		if err := src.Close(); err != nil {
			return err
		}
	}

	return nil
}
