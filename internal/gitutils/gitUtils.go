package gitutils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/jwalton/kitsch/internal/fileutils"
)

// ErrNoGit is emitted when we need to run the git executable, but git is not installed.
var ErrNoGit = errors.New("Git is not installed")

// GitUtils is an object that allows you to retrieve information about
// a git repository.
type gitUtils struct {
	// pathToGit is the path to the git executable.
	pathToGit string
	// The go-git/v5 storer.
	storer *filesystem.Storage
	// fsys is an fs.FS instance bound to the root of the git repository.
	fsys fs.FS
	// RepoRoot is the root folder of the git repository.
	repoRoot string
}

// HeadInfo contains information about the current head.
type HeadInfo struct {
	// HeadDescription is the name of the branch we are currently on if the head
	// is not detached.  If the head is detached, this will be the branch name
	// if we are in the middle of a rebase or merge, the tag name if the head is
	// at a tag, or the short hash otherwise.
	Description string
	// Detached is true if the head is detached.
	Detached bool
	// Hash is the current hash of the head.
	Hash string
	// IsTag is true if the current head matches a tag.
	IsTag bool
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
	// localRef is compared to remoteRef.
	GetAheadBehind(localRef string, remoteRef string) (ahead int, behind int, err error)
	// Head returns information about the current head.
	Head(maxTagsToSearch int) (head HeadInfo, err error)
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

	dotGitPath := filepath.Join(gitRoot, ".git")
	dotGitFs := osfs.New(dotGitPath)
	storer := filesystem.NewStorage(dotGitFs, cache.NewObjectLRUDefault())

	return &gitUtils{
		pathToGit: pathToGit,
		storer:    storer,
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
	config, err := g.storer.Config()
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

// GetAheadBehind returns how many commits ahead and behind the given
// localRef is compared to remoteRef.
func (g *gitUtils) GetAheadBehind(localRef string, remoteRef string) (ahead int, behind int, err error) {
	// If branch and compareToBranch are the same hash, we're done.
	branchRef, err := g.storer.Reference(plumbing.ReferenceName(localRef))
	if err == nil {
		compareBranchRef, err := g.storer.Reference(plumbing.ReferenceName(remoteRef))
		if err == nil {
			if branchRef.Hash() == compareBranchRef.Hash() {
				return 0, 0, nil
			}
		}
	}

	// If not, we need to shell-out to git to find the answer.
	// TODO: Rewrite this as native, as there is quite a bit of overhead going to the shell.
	aheadBehind, err := g.git("rev-list", "--left-right", "--count", localRef+"..."+remoteRef)
	if err != nil {
		return 0, 0, err
	}

	fmt.Sscanf(aheadBehind, "%d %d", &ahead, &behind)
	return ahead, behind, nil
}
