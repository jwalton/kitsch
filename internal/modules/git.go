package modules

import (
	"fmt"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/kitsch-prompt/internal/env"
	"github.com/jwalton/kitsch-prompt/internal/gitutils"
	"github.com/jwalton/kitsch-prompt/internal/style"
)

// GitConfig is configuration for the git module.
type GitConfig struct {
	CommonConfig
}

type gitModule struct {
	config GitConfig
}

// NewGitModule creates a git module.
//
// The default implementation of the git module is loosely based  on
// https://github.com/lyze/posh-git-sh and https://github.com/dahlbyk/posh-git.
//
func NewGitModule(config GitConfig) Module {
	return gitModule{config}
}

// Execute runs a git module.
func (mod gitModule) Execute(env env.Env) ModuleResult {
	config := mod.config

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
		"shortName":  state.Base,
		"stashCount": stashCount,
		"state":      state,
		"stats":      stats,
		"ahead":      ahead,
		"behind":     behind,
		"symbol":     symbol,
	}

	defaultOutput := mod.renderDefault(symbol, state, stats, stashCount, upstream, ahead, behind)

	return executeModule(config.CommonConfig, data, config.Style, defaultOutput)
}

func (mod gitModule) renderDefault(
	symbol string,
	state gitutils.RepositoryState,
	stats gitutils.GitStats,
	stashCount int,
	upstream string,
	ahead int,
	behind int,
) string {
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

func (mod gitModule) renderStats(stats gitutils.GitFileStats) string {
	if stats.Added > 0 || stats.Modified > 0 || stats.Deleted > 0 {
		return fmt.Sprintf("+%d ~%d -%d ", stats.Added, stats.Modified, stats.Deleted)
	}
	return ""
}
