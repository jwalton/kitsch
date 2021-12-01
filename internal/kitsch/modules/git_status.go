package modules

import (
	"fmt"
	"strings"

	"github.com/jwalton/kitsch-prompt/internal/gitutils"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/log"
	"gopkg.in/yaml.v3"
)

// GitStatusModule shows the current status of a git module.
//
// The default implementation of the git status module is loosely based on
// https://github.com/lyze/posh-git-sh and https://github.com/dahlbyk/posh-git.
//
type GitStatusModule struct {
	CommonConfig `yaml:",inline"`
	// IndexStyle is the style to use for the index status.
	IndexStyle string `yaml:"indexStyle"`
	// UnstagedStyle is the style to use for the unstaged file status.
	UnstagedStyle string `yaml:"unstagedStyle"`
	// StashStyle is the style to use for the stash count.
	StashStyle string `yaml:"stashStyle"`
}

type gitStatusModuleResult struct {
	// Index is a `{ Added, Modified, Deleted }` object.  Each is an `int`
	// representing the number of files in the index in that state.
	Index gitutils.GitFileStats
	// Unstaged is a `{ Added, Modified, Deleted }` object.  Each is an `int`
	// representing the number of unstaged files in that state.
	Unstaged gitutils.GitFileStats
	// Unmerged is the total number of unmerged paths in the git repo.
	Unmerged int
	// StashCount is the number of stashes in the git repo.
	StashCount int
}

// Execute runs a git module.
func (mod GitStatusModule) Execute(context *Context) ModuleResult {
	git := context.Git()

	if git == nil {
		return executeModule(context, mod.CommonConfig, gitStatusModuleResult{}, mod.Style, "")
	}

	stats, _ := git.Stats()
	stashCount, err := git.GetStashCount()
	if err != nil {
		stashCount = 0
		log.Warn("Error getting stash count: ", err)
	}

	data := gitStatusModuleResult{
		Index:      stats.Index,
		Unstaged:   stats.Unstaged,
		Unmerged:   stats.Unmerged,
		StashCount: stashCount,
	}

	defaultOutput := mod.renderDefault(context, stats, stashCount)

	return executeModule(context, mod.CommonConfig, data, mod.Style, defaultOutput)
}

func (mod GitStatusModule) renderDefault(
	context *Context,
	stats gitutils.GitStats,
	stashCount int,
) string {
	parts := []string{}
	indexTotal := stats.Index.Added + stats.Index.Modified + stats.Index.Deleted
	unstagedTotal := stats.Unstaged.Added + stats.Unstaged.Modified + stats.Unstaged.Deleted

	indexStyle := defaultStyle(context, mod.IndexStyle, "green")
	unstagedStyle := defaultStyle(context, mod.UnstagedStyle, "red")
	stashStyle := defaultStyle(context, mod.StashStyle, "brightRed")

	if (indexTotal) > 0 || stats.Unmerged > 0 {
		indexPart := mod.renderStats(stats.Index)
		if stats.Unmerged > 0 {
			indexPart += fmt.Sprintf(" !%d", stats.Unmerged)
		}
		indexStats := indexStyle.Apply(indexPart)
		parts = append(parts, indexStats)
	}

	if indexTotal > 0 && unstagedTotal > 0 {
		parts = append(parts, "|")
	}

	if (unstagedTotal) > 0 {
		unstagedStats := unstagedStyle.Apply(mod.renderStats(stats.Unstaged))
		parts = append(parts, unstagedStats)
	}

	if stashCount > 0 {
		stashCountStr := stashStyle.Apply(fmt.Sprintf("(%d)", stashCount))
		parts = append(parts, stashCountStr)
	}

	return strings.Join(parts, " ")
}

func (mod GitStatusModule) renderStats(stats gitutils.GitFileStats) string {
	return fmt.Sprintf("+%d ~%d -%d", stats.Added, stats.Modified, stats.Deleted)
}

func init() {
	registerFactory("git_status", func(node *yaml.Node) (Module, error) {
		var module GitStatusModule
		err := node.Decode(&module)
		return &module, err
	})
}
