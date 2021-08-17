package modules

import (
	"fmt"

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
// • state - A `{ State, Step, Total, Base, Branch }` object.  All of these
//   values are strings.  State is the current state of the repo (e.g. "MERGING"
//   if in the middle of a merge).  Step and Total represent the number of steps
//   left to complete the current operation (e.g. the number of commits left
//   to apply in an interactive rebase), or empty string if no such operation is
//   in progress.  Base is the name of the base branch we are merging from or
//   rebasing from.  Branch is the name of the current branch or hash.
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
		"state":  state,
		"ahead":  ahead,
		"behind": behind,
		"symbol": symbol,
	}

	defaultOutput := mod.renderDefault(symbol, state, upstream, ahead, behind)

	return executeModule(mod.CommonConfig, data, mod.Style, defaultOutput)
}

func (mod GitModule) renderDefault(
	symbol string,
	state gitutils.RepositoryState,
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

	return branch
}

func init() {
	registerFactory("git", func(node *yaml.Node) (Module, error) {
		var module GitModule
		err := node.Decode(&module)
		return &module, err
	})
}
