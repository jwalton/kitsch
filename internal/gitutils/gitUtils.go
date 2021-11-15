package gitutils

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
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
	repo     *git.Repository
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

	repo, err := git.PlainOpen(gitRoot)
	if err != nil {
		return nil
	}

	return &GitUtils{
		pathToGit: pathToGit,
		fsys:      fsys,
		RepoRoot:  gitRoot,
		repo:      repo,
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

// GetCurrentRepo returns a git repo for the current folder, or nil if we are not
// inside a git repo.
func openRepo(path string) *git.Repository {
	gitFolder := fileutils.FindFileInAncestors(path, ".git")

	repo, err := git.PlainOpen(gitFolder)

	if err != nil {
		return nil
	}

	return repo
}
