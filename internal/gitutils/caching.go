package gitutils

// caching is a gitutils that caches results - it assumes the underlying repo
// is not going to change between calls.
type caching struct {
	// TOOD: Make this thread-safe?
	underlying Git

	stashCountInitialized bool
	stashCount            int

	localBranch    string
	upstreamBranch string

	aheadBehindLocalRef  string
	aheadBehindRemoteRef string
	ahead                int
	behind               int

	headInfoTagsSearched int
	headInfo             *HeadInfo
	state                *RepositoryState
	stats                *GitStats
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
	if !c.stashCountInitialized {
		var err error
		c.stashCount, err = c.underlying.GetStashCount()
		if err != nil {
			return 0, err
		}
		c.stashCountInitialized = true
	}
	return c.stashCount, nil
}

// GetUpstream returns the upstream of the current branch if one exists, or
// an empty string otherwise.
func (c *caching) GetUpstream(branch string) string {
	if c.localBranch != branch {
		c.localBranch = branch
		c.upstreamBranch = c.underlying.GetUpstream(branch)
	}
	return c.upstreamBranch
}

// GetAheadBehind returns how many commits ahead and behind the given
// localRef is compared to remoteRef.
func (c *caching) GetAheadBehind(localRef string, remoteRef string) (ahead int, behind int, err error) {
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
	if c.state == nil {
		state := c.underlying.State()
		c.state = &state
	}
	return *c.state
}

// Stats returns status counters for the given git repo.
func (c *caching) Stats() (GitStats, error) {
	if c.stats == nil {
		stats, err := c.underlying.Stats()
		if err != nil {
			return GitStats{}, err
		}
		c.stats = &stats
	}

	return *c.stats, nil
}
