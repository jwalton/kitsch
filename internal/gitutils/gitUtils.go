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
	"strings"

	"github.com/jwalton/kitsch/internal/fileutils"
)

// ErrNoGit is emitted when we need to run the git executable, but git is not installed.
var ErrNoGit = errors.New("Git is not installed")

// GitUtils is an object that allows you to retrieve information about
// a git repository.
type gitUtils struct {
	// pathToGit is the path to the git executable.
	pathToGit string
	// fsys is an fs.FS instance bound to the root of the git repository.
	fsys fs.FS
	// RepoRoot is the root folder of the git repository.
	repoRoot string
}

// Git is an interface for interacting with a git repository.
type Git interface {
	// RepoRoot returns the root of the git repository.
	RepoRoot() string
	// GetStashCount returns the number of stashes.
	GetStashCount() (int, error)
	// GetUpstream returns the upstream of the current branch if one exists, or
	// an empty string otherwise.
	GetUpstream(branch string) string
	// GetAheadBehind returns how many commits ahead and behind the given
	// branch is compared to compareToBranch.  You can use `HEAD` for the branch name.
	GetAheadBehind(branch string, compareToBranch string) (ahead int, behind int, err error)
	// State returns the current state of the repository.
	State() RepositoryState
	// Stats returns status counters for the given git repo.
	Stats() (GitStats, error)
}

// New returns a new instance of `GitUtils` for the specified folder.
// If the folder is not a git repository, it will return nil.
func New(pathToGit string, folder string) Git {
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

	return &gitUtils{
		pathToGit: pathToGit,
		fsys:      fsys,
		repoRoot:  gitRoot,
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
func (g *gitUtils) git(args ...string) (string, error) {
	if g.pathToGit == "" {
		return "", ErrNoGit
	}

	cmd := exec.Command(g.pathToGit, args...)
	cmd.Dir = g.repoRoot

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (g *gitUtils) RepoRoot() string {
	return g.repoRoot
}

func (g *gitUtils) GetStashCount() (int, error) {
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

func (g *gitUtils) GetUpstream(branch string) string {
	config, err := g.localConfig()
	if err != nil {
		return ""
	}

	// TODO: If `branch` is HEAD, resolve it.
	branchConfig := config.Branches[branch]
	if branchConfig == nil {
		return ""
	}

	if !strings.HasPrefix(branchConfig.Merge, "refs/heads/") {
		return ""
	}

	return branchConfig.Remote + "/" + branchConfig.Merge[11:]
}

func (g *gitUtils) GetAheadBehind(branch string, compareToBranch string) (ahead int, behind int, err error) {
	aheadBehind, err := g.git("rev-list", "--left-right", "--count", branch+"..."+compareToBranch)

	if err != nil {
		return 0, 0, err
	}

	fmt.Sscanf(aheadBehind, "%d %d", &ahead, &behind)
	return ahead, behind, nil
}

// ReadObject reads a git object from the repo.
func (g *gitUtils) ReadObject(hash string) (objectType string, data []byte, err error) {
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
func (g *gitUtils) ReadObjectOfType(objectType string, hash string) (data []byte, err error) {
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

// getContent reads a "length-of-content", then a null byte, then the content.
func getContent(from []byte) ([]byte, error) {
	// Read the size of the contents.
	sizeStr, sizeLen := readUntil(from, '\x00')
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid size in git object: %s, %w", sizeStr, err)
	}

	content := from[sizeLen+1:]
	if int64(len(content)) != size {
		return nil, fmt.Errorf("expected %d bytes but only got %d", size, len(content))
	}

	return content, nil
}

// Read bytes from "bytes" until we find the specified byte.  Returns the read data
// as a string up until but not including the `until` byte.
func readUntil(bytes []byte, until byte) (string, int) {
	index := 0
	for index < len(bytes) && bytes[index] != until {
		index++
	}
	return string(bytes[:index]), index
}
