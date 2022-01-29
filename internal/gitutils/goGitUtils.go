package gitutils

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// This function was copied from github.com/go-git/go-git/v5 so we don't have
// to depend on go-git directly.
func (g *gitUtils) tags() (storer.ReferenceIter, error) {
	refIter, err := g.storer.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(
		func(r *plumbing.Reference) bool {
			return r.Name().IsTag()
		}, refIter), nil
}
