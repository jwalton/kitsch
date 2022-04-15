package gitutils

import "sync"

// caching is a gitutils that caches results - it assumes the underlying repo
// is not going to change between calls.
type caching struct {
	underlying Git

	mutex sync.Mutex

	stashCountOnce sync.Once
	stashCount     int
	stashCountErr  error

	localBranch    string
	upstreamBranch string

	aheadBehindLocalRef  string
	aheadBehindRemoteRef string
	ahead                int
	behind               int

	headInfoTagsSearched int
	headInfo             *HeadInfo
	stateOnce            sync.Once
	state                *RepositoryState
	statsOnce            sync.Once
	stats                GitStats
	statsError           error
}

// NewCaching returns a new caching instance of Git.  The returned instance
// assumes the repo does not change between calls, so will not recompute the
// same values more than once.
func NewCaching(pathToGit string, folder string) Git {
	underlying := New(pathToGit, folder)
	if underlying == nil {
		return nil
	}
	return &caching{underlying: underlying}
}

// RepoRoot returns the root of the git repository.
func (c *caching) RepoRoot() string {
	return c.underlying.RepoRoot()
}

// GetStashCount returns the number of stashes.
func (c *caching) GetStashCount() (int, error) {
	c.stashCountOnce.Do(func() {
		c.stashCount, c.stashCountErr = c.underlying.GetStashCount()
	})
	return c.stashCount, c.stashCountErr
}

// GetUpstream returns the upstream of the current branch if one exists, or
// an empty string otherwise.
func (c *caching) GetUpstream(branch string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.localBranch != branch {
		c.localBranch = branch
		c.upstreamBranch = c.underlying.GetUpstream(branch)
	}
	return c.upstreamBranch
}

// GetAheadBehind returns how many commits ahead and behind the given
// localRef is compared to remoteRef.
func (c *caching) GetAheadBehind(localRef string, remoteRef string) (ahead int, behind int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.aheadBehindLocalRef != localRef || c.aheadBehindRemoteRef != remoteRef {
		ahead, behind, err = c.underlying.GetAheadBehind(localRef, remoteRef)
		if err != nil {
			return 0, 0, err
		}
		c.aheadBehindLocalRef = localRef
		c.aheadBehindRemoteRef = remoteRef
		c.ahead = ahead
		c.behind = behind
	}
	return c.ahead, c.behind, nil
}

// Head returns information about the current head.
func (c *caching) Head(maxTagsToSearch int) (head HeadInfo, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	haveHeadInfo := c.headInfo != nil && (!c.headInfo.Detached || c.headInfo.IsTag || maxTagsToSearch < c.headInfoTagsSearched)
	if !haveHeadInfo {
		c.headInfoTagsSearched = maxTagsToSearch
		headInfo, err := c.underlying.Head(maxTagsToSearch)
		if err != nil {
			return HeadInfo{}, err
		}
		c.headInfo = &headInfo
	}
	return *c.headInfo, nil
}

// State returns the current state of the repository.
func (c *caching) State() RepositoryState {
	c.stateOnce.Do(func() {
		state := c.underlying.State()
		c.state = &state
	})
	return *c.state
}

// Stats returns status counters for the given git repo.
func (c *caching) Stats() (GitStats, error) {
	c.statsOnce.Do(func() {
		c.stats, c.statsError = c.underlying.Stats()
	})
	return c.stats, c.statsError
}
