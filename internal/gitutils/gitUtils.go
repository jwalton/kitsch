package gitutils

import (
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
)

// ErrNoGit is emitted when we need to run the git executable, but git is not installed.
var ErrNoGit = errors.New("Git is not installed")

// GitUtils is an object that allows you to retrieve information about
// a git repository.
type GitUtils struct {
	// pathToGit is the path to the git executable.
	pathToGit string
	// fsys is an fs.FS instance bound to the root of the git repository.
	fsys fs.FS
	// RepoRoot is the root folder of the git repository.
	RepoRoot string
}

// New returns a new instance of `GitUtils` for the specified repository.
func New(pathToGit string, folder string) *GitUtils {
	// Resolve the path to the git executable
	pathToGit, err := fileutils.LookPathSafe(pathToGit)
	if err != nil {
		pathToGit = ""
	}

	// Figure out whether or not we're inside a git repo.
	gitRoot := FindGitRoot(folder)

	var fsys fs.FS = nil
	if gitRoot != "" {
		fsys = os.DirFS(gitRoot)
	}

	if gitRoot == "" {
		return nil
	}

	return &GitUtils{
		pathToGit: pathToGit,
		fsys:      fsys,
		RepoRoot:  gitRoot,
	}
}

// FindGitRoot returns the root of the current git repo.
func FindGitRoot(cwd string) string {
	gitFolder := fileutils.FindFileInAncestors(cwd, ".git")
	if gitFolder != "" {
		return filepath.Dir(gitFolder)
	}
	return ""
}

// git will run a git command in the root folder of the git repository.
// Returns empty string if there was an error running the command.
func (g *GitUtils) git(args ...string) (string, error) {
	if g.pathToGit == "" {
		return "", ErrNoGit
	}

	cmd := exec.Command(g.pathToGit, args...)
	cmd.Dir = g.RepoRoot

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func countLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0

	for {
		n, err := r.Read(buf)
		for i := 0; i < n; i++ {
			if buf[i] == '\n' {
				count++
			}
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}
	}

	return count, nil
}

// GetStashCount returns the number of stashes, or 0 if there are none or
// the path is not a git repo.
//
// `path` should be the git root folder.
func (g *GitUtils) GetStashCount() (int, error) {
	// TODO: Read .git/logs/refs/stash, and count the number of `\n`s.`
	file, err := g.fsys.Open(".git/logs/refs/stash")
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	defer file.Close()
	return countLines(file)
}

// ReadObject reads a git object from the repo.
func (g *GitUtils) ReadObject(hash string) (objectType string, data []byte, err error) {
	filename := ".git/objects/" + hash[0:2] + "/" + hash[2:]
	file, err := g.fsys.Open(filename)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return "", nil, err
	}
	defer zlibReader.Close()

	content, err := io.ReadAll(zlibReader)
	if err != nil {
		return "", nil, err
	}

	// Read the object type.
	objectType, objectTypeLen := readUntil(content, ' ')

	// Read the size of the contents.
	data, err = getContent(content[objectTypeLen+1:])
	if err != nil {
		return "", nil, err
	}

	return objectType, data, nil
}

// ErrIncorrectType is returned when an object is not of the requested type.
var ErrIncorrectType = errors.New("incorrect object type")

// ReadObjectOfType reads a git object from repo, but only if it is of the
// provided type.  If the object is not the correct type, returns ErrIncorrectType.
func (g *GitUtils) ReadObjectOfType(objectType string, hash string) (data []byte, err error) {
	filename := ".git/objects/" + hash[0:2] + "/" + hash[2:]
	file, err := g.fsys.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	// Read the file type.
	fileType := make([]byte, len(objectType)+1)
	n, err := zlibReader.Read(fileType)
	if n != len(objectType)+1 || string(fileType[0:len(objectType)]) != objectType {
		return nil, ErrIncorrectType
	}

	content, err := io.ReadAll(zlibReader)

	// Read the size of the contents.
	return getContent(content)
}

// getContent reads a length-of-content, then a null byte, then the content.
func getContent(content []byte) ([]byte, error) {
	// Read the size of the contents.
	sizeStr, sizeLen := readUntil(content, '\x00')
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid size in git object: %s, %w", sizeStr, err)
	}

	data := content[sizeLen+1:]
	if len(data) != int(size) {
		return nil, fmt.Errorf("expected %d bytes but only got %d", size, len(data))
	}

	return data, nil
}

func readUntil(bytes []byte, until byte) (string, int) {
	index := 0
	for index < len(bytes) && bytes[index] != until {
		index++
	}
	return string(bytes[:index]), index
}
