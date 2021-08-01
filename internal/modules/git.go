package modules

import (
	"fmt"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/gitutils"
	"github.com/jwalton/kitsch-prompt/internal/style"
	"gopkg.in/yaml.v3"
)

// GitModule shows information about the current git repo.
//
// The default implementation of the git module is loosely based  on
// https://github.com/lyze/posh-git-sh and https://github.com/dahlbyk/posh-git.
//
// Provides the following template variables:
//
// • stashCount - The number of stashes.
//
// • state - A `{ State, Step, Total, Base, Branch }` object.  All of these
//   values are strings.  State is the current state of the repo (e.g. "MERGING"
//   if in the middle of a merge).  Step and Total represent the number of steps
//   left to complete the current operation (e.g. the number of commits left
//   to apply in an interactive rebase), or empty string if no such operation is
//   in progress.  Base is the name of the base branch we are merging from or
//   rebasing from.  Branch is the name of the current branch or hash.
//
// • stats - A `{ Index, Files, Unmerged }` object.  Index and Files are the
//   number of `{ Added, Modified, Deleted }` files in the index and working
//   directory, respectively.  Unmerged is the number of unmerged files, if
//   there is a merge operation in progress.
//
// • ahead - The number of commits ahead of the upstream branch.
//
// • behind - The number of commits behind the upstream branch.
//
// • symbol - The symbol to use to indicate the current state of the repo.
//
type GitModule struct {
	CommonConfig `yaml:",inline"`
}

// Execute runs a git module.
func (mod GitModule) Execute(env env.Env) ModuleResult {
	cwd := env.Getwd()
	git := gitutils.New("git", cwd)

	if git == nil {
		return ModuleResult{}
	}

	state := git.State()
	var ahead, behind int
	var upstream string

	if state.State == gitutils.StateNone && state.Branch != "" {
		upstream = git.GetUpstream(state.Branch)
		if upstream != "" {
			ahead, behind, _ = git.GetAheadBehind("HEAD", upstream)
		}
	}

	stats, _ := git.Stats()

	stashCount := git.GetStashCount()

	symbol := "?"
	if upstream != "" {
		if ahead > 0 && behind > 0 {
			symbol = "↕"
		} else if ahead > 0 {
			symbol = "↑"
		} else if behind > 0 {
			symbol = "↓"
		} else {
			symbol = "≡"
		}
	}

	data := map[string]interface{}{
		"stashCount": stashCount,
		"state":      state,
		"stats":      stats,
		"ahead":      ahead,
		"behind":     behind,
		"symbol":     symbol,
	}

	defaultOutput := mod.renderDefault(symbol, state, stats, stashCount, upstream, ahead, behind)

	return executeModule(mod.CommonConfig, data, mod.Style, defaultOutput)
}

func (mod GitModule) renderDefault(
	symbol string,
	state gitutils.RepositoryState,
	stats gitutils.GitStats,
	stashCount int,
	upstream string,
	ahead int,
	behind int,
) string {
	// TODO: Make these styles configurable.
	branchStyle := style.Style{FG: "brightCyan"}
	if upstream != "" {
		if ahead > 0 && behind > 0 {
			branchStyle = style.Style{FG: "brightYellow"}
		} else if ahead > 0 {
			branchStyle = style.Style{FG: "brightGreen"}
		} else if behind > 0 {
			branchStyle = style.Style{FG: "brightRed"}
		}
	}

	branch := state.Base + " " + symbol
	if state.State != gitutils.StateNone {
		branch = branch + "|" + string(state.State)
		if state.Total != "" {
			branch = fmt.Sprintf("%s %s/%s", branch, state.Step, state.Total)
		}
	}
	branch, _, _, _ = branchStyle.Apply(branch)

	indexStats := gchalk.Green(mod.renderStats(stats.Index))
	fileStats := gchalk.Red(mod.renderStats(stats.Files))

	statsJoiner := ""
	if fileStats != "" && indexStats != "" {
		statsJoiner = "| "
	}

	unmergedStats := ""
	if stats.Unmerged > 0 {
		unmergedStats = gchalk.BrightMagenta(fmt.Sprintf("%d! ", stats.Unmerged))
	}

	stashCountStr := ""
	if stashCount > 0 {
		stashCountStr = gchalk.BrightRed(fmt.Sprintf("(%d)", stashCount))
	}

	statsPart := fmt.Sprintf("%s%s%s%s%s", indexStats, statsJoiner, fileStats, unmergedStats, stashCountStr)
	if statsPart != "" {
		statsPart = " " + statsPart
	}

	return "[" + branch + statsPart + "]"
}

func (mod GitModule) renderStats(stats gitutils.GitFileStats) string {
	if stats.Added > 0 || stats.Modified > 0 || stats.Deleted > 0 {
		return fmt.Sprintf("+%d ~%d -%d ", stats.Added, stats.Modified, stats.Deleted)
	}
	return ""
}

func init() {
	registerFactory("git", func(node *yaml.Node) (Module, error) {
		var module GitModule
		err := node.Decode(&module)
		return &module, err
	})
}
