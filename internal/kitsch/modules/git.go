package modules

import (
	"fmt"
	"strings"

	"github.com/jwalton/kitsch/internal/gitutils"
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas GitModule

// GitModule shows information about the current git repo.
//
// The default implementation of the git module is loosely based  on
// https://github.com/lyze/posh-git-sh and https://github.com/dahlbyk/posh-git.
//
type GitModule struct {
	CommonConfig `yaml:",inline"`
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",enum=git"`
}

type gitResult struct {
	// State is the current state of the repo.
	State gitutils.RepositoryState `json:"state"`
	// Upstream is the name of the upstream branch (e.g. "origin/master"), or ""
	// if there is no upstream or we are currently detached.
	Upstream string `json:"upstream"`
	// Ahead is the number of commits we are ahead of the upstream branch, or 0 if there is no upstream branch.
	Ahead int `json:"ahead"`
	// Behind is the number of commits we are behind of the upstream branch, or 0 if there is no upstream branch.
	Behind int `json:"behind"`
	// Symbol is the symbol to use to indicate the current state of the repo.
	Symbol string `json:"symbol"`
	// AheadBehind is "ahead" if we are ahead of the upstream branch, "behind"
	// if we are behind, "diverged" if we are both, and "" otherwise.
	AheadBehind string `json:"aheadBehind"`
}

// Execute runs a git module.
func (mod GitModule) Execute(context *Context) ModuleResult {
	git := context.Git()

	if git == nil {
		return ModuleResult{}
	}

	state := git.State()
	var ahead, behind int
	var upstream string

	if !state.IsDetached {
		upstream = git.GetUpstream(state.HeadDescription)
		if upstream != "" {
			ahead, behind, _ = git.GetAheadBehind("refs/heads/"+state.HeadDescription, "refs/remotes/"+upstream)
		}
	}

	symbol := "?"
	aheadBehind := ""
	if upstream != "" {
		if ahead > 0 && behind > 0 {
			symbol = "↕"
			aheadBehind = "diverged"
		} else if ahead > 0 {
			symbol = "↑"
			aheadBehind = "ahead"
		} else if behind > 0 {
			symbol = "↓"
			aheadBehind = "behind"
		} else {
			symbol = "≡"
		}
	}

	data := gitResult{
		State:       state,
		Upstream:    upstream,
		Ahead:       ahead,
		Behind:      behind,
		Symbol:      symbol,
		AheadBehind: aheadBehind,
	}

	defaultOutput := mod.renderDefault(context, symbol, state, upstream, ahead, behind)

	return executeModule(context, mod.CommonConfig, data, mod.Style, defaultOutput)
}

func (mod GitModule) renderDefault(
	context *Context,
	symbol string,
	state gitutils.RepositoryState,
	upstream string,
	ahead int,
	behind int,
) string {
	out := strings.Builder{}

	out.WriteString(state.HeadDescription)

	if behind > 0 {
		out.WriteString(fmt.Sprintf(" ↓%d", behind))
	}
	if ahead > 0 {
		out.WriteString(fmt.Sprintf(" ↑%d", ahead))
	}
	if behind == 0 && ahead == 0 {
		out.WriteString(" " + symbol)
	}

	if state.State != gitutils.StateNone {
		out.WriteString("|" + string(state.State))
		if state.Total != "" {
			out.WriteString(fmt.Sprintf(" %s/%s", state.Step, state.Total))
		}
	}

	return out.String()
}

func init() {
	registerModule(
		"git",
		registeredModule{
			jsonSchema: schemas.GitModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := GitModule{Type: "git"}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
