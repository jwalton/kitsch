package modules

import (
	"fmt"
	"strings"

	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/gitutils"
	"github.com/jwalton/kitsch-prompt/internal/style"
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
	IndexStyle style.Style `yaml:"indexStyle"`
	// UnstagedStyle is the style to use for the unstaged file status.
	UnstagedStyle style.Style `yaml:"unstagedStyle"`
	// UnmergedStyle is the style to use for the unmerged files count.
	UnmergedStyle style.Style `yaml:"unmergedStyle"`
	// StashStyle is the style to use for the stash count.
	StashStyle style.Style `yaml:"stashStyle"`
}

// Execute runs a git module.
func (mod GitStatusModule) Execute(env env.Env) ModuleResult {
	git := env.Git()

	if git == nil {
		return ModuleResult{}
	}

	stats, _ := git.Stats()
	stashCount := git.GetStashCount()

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

	defaultOutput := mod.renderDefault(stats, stashCount)

	return executeModule(mod.CommonConfig, data, mod.Style, defaultOutput)
}

func (mod GitStatusModule) renderDefault(
	stats gitutils.GitStats,
	stashCount int,
) string {
	parts := []string{}
	indexTotal := stats.Index.Added + stats.Index.Modified + stats.Index.Deleted
	filesTotal := stats.Unstaged.Added + stats.Unstaged.Modified + stats.Unstaged.Deleted

	if (indexTotal) > 0 {
		indexStats, _, _, _ := mod.IndexStyle.Default(style.Style{FG: "green"}).Apply(mod.renderStats(stats.Index))
		parts = append(parts, indexStats)
	}

	if indexTotal > 0 && filesTotal > 0 {
		parts = append(parts, "|")
	}

	if (filesTotal) > 0 {
		fileStats, _, _, _ := mod.UnstagedStyle.Default(style.Style{FG: "red"}).Apply(mod.renderStats(stats.Unstaged))
		parts = append(parts, fileStats)
	}

	if stats.Unmerged > 0 {
		unmergedStats, _, _, _ := mod.UnmergedStyle.Default(style.Style{FG: "brightMagenta"}).Apply(fmt.Sprintf("%d!", stats.Unmerged))
		parts = append(parts, unmergedStats)
	}

	if stashCount > 0 {
		stashCountStr, _, _, _ := mod.StashStyle.Default(style.Style{FG: "brightRed"}).Apply(fmt.Sprintf("(%d)", stashCount))
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
