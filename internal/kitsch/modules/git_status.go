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
// Configuration:
//
// • unstagedStyle - The style to use for the unstaged file status.  Defaults to red.
//
// • indexStyle - The style to use for the index status.  Defaults to green.
//
// • unmergedStyle - The style to use for the unmerged files count.  Defaults to bright magenta.
//
// • stashStyle - The style to use for the stasg count.  Defaults to bright red.
//
// Provides the following template variables:
//
// • Index - An `{ Added, Modified, Deleted, Total }` object.  Each is an `int`
//   representing the number of files in the index in that state.
//
// • Unstaged - An `{ Added, Modified, Deleted, Total }` object.  Each is an `int`
//   representing the number of unstaged files in that state.
//
// • Unmerged - An `int` representing the number of unmerged files.
//
// • StashCount - An `int` representing the number of stashes.
//
type GitStatusModule struct {
	CommonConfig `yaml:",inline"`
	// IndexStyle is the style to use for the index status.
	IndexStyle string `yaml:"indexStyle"`
	// UnstagedStyle is the style to use for the unstaged file status.
	UnstagedStyle string `yaml:"unstagedStyle"`
	// UnmergedStyle is the style to use for the unmerged files count.
	UnmergedStyle string `yaml:"unmergedStyle"`
	// StashStyle is the style to use for the stash count.
	StashStyle string `yaml:"stashStyle"`
}

// Execute runs a git module.
func (mod GitStatusModule) Execute(context *Context) ModuleResult {
	git := context.Environment.Git()

	if git == nil {
		return ModuleResult{}
	}

	stats, _ := git.Stats()
	stashCount, err := git.GetStashCount()
	if err != nil {
		stashCount = 0
		log.Warn("Error getting stash count: ", err)
	}

	data := map[string]interface{}{
		"Index": map[string]interface{}{
			"Added":    stats.Index.Added,
			"Modified": stats.Index.Modified,
			"Deleted":  stats.Index.Deleted,
			"Total":    stats.Index.Added + stats.Index.Modified + stats.Index.Deleted,
		},
		"Unstaged": map[string]interface{}{
			"Added":    stats.Unstaged.Added,
			"Modified": stats.Unstaged.Modified,
			"Deleted":  stats.Unstaged.Deleted,
			"Total":    stats.Unstaged.Added + stats.Unstaged.Modified + stats.Unstaged.Deleted,
		},
		"Unmerged":   stats.Unmerged,
		"StashCount": stashCount,
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
	filesTotal := stats.Unstaged.Added + stats.Unstaged.Modified + stats.Unstaged.Deleted

	indexStyle := defaultStyle(context, mod.IndexStyle, "green")
	unstagedStyle := defaultStyle(context, mod.UnstagedStyle, "red")
	unmergedStyle := defaultStyle(context, mod.UnmergedStyle, "brightMagenta")
	stashStyle := defaultStyle(context, mod.StashStyle, "brightRed")

	if (indexTotal) > 0 {
		indexStats := indexStyle.Apply(mod.renderStats(stats.Index))
		parts = append(parts, indexStats)
	}

	if indexTotal > 0 && filesTotal > 0 {
		parts = append(parts, "|")
	}

	if (filesTotal) > 0 {
		fileStats := unstagedStyle.Apply(mod.renderStats(stats.Unstaged))
		parts = append(parts, fileStats)
	}

	if stats.Unmerged > 0 {
		unmergedStats := unmergedStyle.Apply(fmt.Sprintf("!%d", stats.Unmerged))
		parts = append(parts, unmergedStats)
	}

	if stashCount > 0 {
		stashCountStr := stashStyle.Apply(fmt.Sprintf("(%d)", stashCount))
		parts = append(parts, stashCountStr)
	}

	return strings.Join(parts, " ")
}

func (mod GitStatusModule) renderStats(stats gitutils.GitFileStats) string {
	if stats.Added > 0 || stats.Modified > 0 || stats.Deleted > 0 {
		return fmt.Sprintf("+%d ~%d -%d", stats.Added, stats.Modified, stats.Deleted)
	}
	return ""
}

func init() {
	registerFactory("git_status", func(node *yaml.Node) (Module, error) {
		var module GitStatusModule
		err := node.Decode(&module)
		return &module, err
	})
}
