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
	// MaxTagsToSearch is the maximum number of tags to search when checking to
	// see if HEAD is a tagged release.  Defaults to 200.
	MaxTagsToSearch int `yaml:"maxTagsToSearch"`
}

type gitResult struct {
	// State is the current state of this repo.
	State gitutils.RepositoryStateType `yaml:"state"`
	// Step is the current step number if we are rebasing, 0 otherwise.
	Step string `yaml:"step"`
	// Total is the total number of steps to complete to finish the rebase, or 0
	// if not rebasing.
	Total string `yaml:"total"`
	// Head is information about the current HEAD.
	Head gitutils.HeadInfo `json:"head"`
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
	// if we are behind, "diverged" if we are both, and "upToDate" otherwise.
	AheadBehind string `json:"aheadBehind"`
}

// Execute runs a git module.
func (mod GitModule) Execute(context *Context) ModuleResult {
	git := context.Git()

	if git == nil {
		return ModuleResult{}
	}

	head, err := git.Head(mod.MaxTagsToSearch)
	if err != nil {
		head.Description = "???"
		head.Detached = true
		head.Hash = "???"
	}

	state := git.State()
	var ahead, behind int
	var upstream string

	if !head.Detached {
		upstream = git.GetUpstream(head.Description)
		if upstream != "" {
			ahead, behind, _ = git.GetAheadBehind("refs/heads/"+head.Description, "refs/remotes/"+upstream)
		}
	}

	symbol := "?"
	aheadBehind := "upToDate"
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
			aheadBehind = "upToDate"
		}
	}

	data := gitResult{
		State:       state.State,
		Step:        state.Step,
		Total:       state.Total,
		Head:        head,
		Upstream:    upstream,
		Ahead:       ahead,
		Behind:      behind,
		Symbol:      symbol,
		AheadBehind: aheadBehind,
	}

	defaultOutput := mod.renderDefault(context, symbol, data)

	return executeModule(context, mod.CommonConfig, data, mod.Style, defaultOutput)
}

func (mod GitModule) renderDefault(
	context *Context,
	symbol string,
	data gitResult,
) string {
	out := strings.Builder{}

	out.WriteString(data.Head.Description)

	if data.Behind > 0 {
		out.WriteString(fmt.Sprintf(" ↓%d", data.Behind))
	}
	if data.Ahead > 0 {
		out.WriteString(fmt.Sprintf(" ↑%d", data.Ahead))
	}
	if data.Behind == 0 && data.Ahead == 0 {
		out.WriteString(" " + symbol)
	}

	if data.State != gitutils.StateNone {
		out.WriteString("|" + string(data.State))
		if data.Total != "" {
			out.WriteString(fmt.Sprintf(" %s/%s", data.Step, data.Total))
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
				module := GitModule{
					Type:            "git",
					MaxTagsToSearch: 200,
				}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
