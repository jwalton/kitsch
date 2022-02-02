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
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=git"`
	// MaxTagsToSearch is the maximum number of tags to search when checking to
	// see if HEAD is a tagged release.  Defaults to 200.
	MaxTagsToSearch int `yaml:"maxTagsToSearch"`

	// RebasingInteractive is a description to show when an interactive rebase in in progress.
	RebasingInteractive string `yaml:"rebaseInteractive"`
	// RebaseMerging is a description to show when a merge in in progress.
	RebaseMerging string `yaml:"rebaseMerging"`
	// Rebasing is a description to show when a rebase operation in in progress.
	Rebasing string `yaml:"rebasing"`
	// AMing is a description to show when an `am` operation in in progress.
	AMing string `yaml:"aming"`
	// RebaseAMing is a description to show when an ambiguous apply-mailbox or rebase is in progress.
	RebaseAMing string `yaml:"rebaseAMing"`
	// Merging is a description to show when a merge in in progress.
	Merging string `yaml:"merging"`
	// CherryPicking is a description to show when a cherry-pick in in progress.
	CherryPicking string `yaml:"cherryPicking"`
	// Reverting is a description to show when a revert in in progress.
	Reverting string `yaml:"reverting"`
	// Bisecting is a description to show when a bisect in in progress.
	Bisecting string `yaml:"bisecting"`
}

type gitResult struct {
	// State is the current state of this repo.
	State string `yaml:"state"`
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
		return ModuleResult{DefaultText: "", Data: gitResult{}}
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
		State:       mod.getStateDescription(state.State),
		Step:        state.Step,
		Total:       state.Total,
		Head:        head,
		Upstream:    upstream,
		Ahead:       ahead,
		Behind:      behind,
		Symbol:      symbol,
		AheadBehind: aheadBehind,
	}

	return ModuleResult{
		DefaultText: mod.renderDefault(context, symbol, data),
		Data:        data,
	}
}

func (mod GitModule) getStateDescription(state gitutils.RepositoryStateType) string {
	switch state {
	case gitutils.StateRebasingInteractive:
		return mod.RebasingInteractive
	case gitutils.StateRebaseMerging:
		return mod.RebaseMerging
	case gitutils.StateRebasing:
		return mod.Rebasing
	case gitutils.StateAMing:
		return mod.AMing
	case gitutils.StateRebaseAMing:
		return mod.RebaseAMing
	case gitutils.StateMerging:
		return mod.Merging
	case gitutils.StateCherryPicking:
		return mod.CherryPicking
	case gitutils.StateReverting:
		return mod.Reverting
	case gitutils.StateBisecting:
		return mod.Bisecting
	default:
		return string(state)
	}
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

	if data.State != "" {
		out.WriteString("|" + data.State)
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
					Type:                "git",
					MaxTagsToSearch:     200,
					RebasingInteractive: "REBASE-i",
					RebaseMerging:       "REBASE-m",
					Rebasing:            "REBASE",
					AMing:               "AM",
					RebaseAMing:         "REBASE/AM",
					Merging:             "MERGING",
					CherryPicking:       "CHERRY-PICKING",
					Reverting:           "REVERTING",
					Bisecting:           "BISECTING",
				}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
