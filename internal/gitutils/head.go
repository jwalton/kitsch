package gitutils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// Head returns information about the current head.
func (g *gitUtils) Head(maxTagsToSearch int) (HeadInfo, error) {
	if g.fsys == nil {
		return HeadInfo{}, fmt.Errorf("no git repo found")
	}

	head, err := storer.ResolveReference(g.storer, plumbing.HEAD)
	if err != nil {
		return HeadInfo{}, err
	}

	headRef := head.Name().String()
	headHash := head.Hash().String()

	isDetached := true
	description := extractBranchName(g.readFileIfExist(".git/rebase-merge/head-name"))
	if description == "" && strings.HasPrefix(headRef, "refs/heads/") {
		// If the HEAD file is a symbolic reference to a branch, extract the branch name.
		description = headRef[11:]
		isDetached = false
	}

	// If we don't have a description, try to get a tag name
	if description == "" {
		tag, err := g.GetTagNameForHash(headHash, maxTagsToSearch)
		if err == nil && tag != "" {
			description = "(" + strings.TrimSpace(tag) + ")"
		}
	}

	// If that fails, use the hash
	if description == "" {
		description = "(" + headHash[0:shortSHALength] + "...)"
	}

	if description == "" {
		description = "???"
	}

	return HeadInfo{
		Description: description,
		Detached:    isDetached,
		Hash:        headHash,
	}, nil
}

// GetTagNameForHash returns the tag name for the hash, or an error if no such
// tag exists.  "hash" can be a short hash.
//
// maxTagsToSearch is the maximum number of tag refs to examine when looking for
// the current hash.  If this is negative, we will search all refs.
func (g *gitUtils) GetTagNameForHash(hash string, maxTagsToSearch int) (string, error) {
	// Check lightweight tags
	tags, err := g.tags()
	if err == nil {
		count := 0
		for ref, err := tags.Next(); err == nil && ref != nil; ref, err = tags.Next() {
			count++
			if maxTagsToSearch >= 0 && count > maxTagsToSearch {
				break
			}
			if g.hashMatchesTag(hash, ref.Hash()) {
				result := strings.TrimPrefix(string(ref.Name()), "refs/tags/")
				return result, nil
			}

		}
	}

	return "", errNotFound
}

// hashMatchesTag returns true if the `hash` represents the same commit as
// the given tag.  `hash` can be a "short hash".  This will return true
// if hash is a prefix of tagHash (the "lightweight" tag case) or if
// the tagHash is an annotated hash and hash is the hash of the tag.
func (g *gitUtils) hashMatchesTag(hash string, tagHash plumbing.Hash) bool {
	// If the tag is a lightweight tag, the hash should match.
	if strings.HasPrefix(tagHash.String(), hash) {
		return true
	}

	// Check to see if this is an annotated tag
	tagObj, err := object.GetTag(g.storer, tagHash)
	if err == nil {
		if strings.HasPrefix(tagObj.Target.String(), hash) {
			return true
		}
	}

	return false
}

var errNotFound = errors.New("Not found")

// extractBranchName returns the branch name from a symbolic ref, or returns
// the passed in string otherwise.
func extractBranchName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.TrimPrefix(name, "refs/heads/")
	return name
}
