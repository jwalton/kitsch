package gitutils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/jwalton/kitsch/internal/fileutils"
)

// ErrNoGit is emitted when we need to run the git executable, but git is not installed.
var ErrNoGit = errors.New("Git is not installed")

// GitUtils is an object that allows you to retrieve information about
// a git repository.
type gitUtils struct {
	// pathToGit is the path to the git executable.
	pathToGit string
	// The git repository.
	repo *git.Repository
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

	repo, err := git.PlainOpen(gitRoot)
	if err != nil {
		return nil
	}

	return &gitUtils{
		pathToGit: pathToGit,
		repo:      repo,
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
	config, err := g.repo.Config()
	if err != nil {
		return ""
	}

	branchConfig := config.Branches[branch]
	if branchConfig == nil {
		return ""
	}

	if !branchConfig.Merge.IsBranch() {
		return ""
	}

	return branchConfig.Remote + "/" + branchConfig.Merge.String()[11:]
}

func (g *gitUtils) GetAheadBehind(branch string, compareToBranch string) (ahead int, behind int, err error) {
	aheadBehind, err := g.git("rev-list", "--left-right", "--count", branch+"..."+compareToBranch)

	if err != nil {
		return 0, 0, err
	}

	fmt.Sscanf(aheadBehind, "%d %d", &ahead, &behind)
	return ahead, behind, nil
}
