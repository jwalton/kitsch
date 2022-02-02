package modules

import (
	"fmt"
	"strings"

	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas GitDiverged

// GitDiverged shows information about whether the current git repo is ahead, behind,
// or has diverged from its upstream.
//
type GitDiverged struct {
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=git_diverged"`
	// AheadSymbol is the symbol to use when the current branch is ahead of its upstream.
	AheadSymbol string `yaml:"aheadSymbol"`
	// BehindSymbol is the symbol to use when the current branch is behind its upstream.
	BehindSymbol string `yaml:"behindSymbol"`
	// DivergedSymbol is the symbol to use when the current branch has diverged from its upstream.
	DivergedSymbol string `yaml:"divergedSymbol"`
	// UpToDateSymbol is the symbol to use when the current branch is up to date with its upstream.
	UpToDateSymbol string `yaml:"upToDateSymbol"`
	// NoUpstreamSymbol is the symbol to use when the current branch has no upstream.
	NoUpstreamSymbol string `yaml:"noUpstreamSymbol"`
}

type gitDivergedResult struct {
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
func (mod GitDiverged) Execute(context *Context) ModuleResult {
	git := context.Git()

	result := gitDivergedResult{}

	if git == nil {
		return ModuleResult{DefaultText: result.Symbol, Data: result}
	}

	head, err := git.Head(0)
	if err != nil {
		return ModuleResult{DefaultText: result.Symbol, Data: result}
	}

	var ahead, behind int
	var upstream string
	if !head.Detached {
		upstream = git.GetUpstream(head.Description)
		if upstream != "" {
			ahead, behind, _ = git.GetAheadBehind("refs/heads/"+head.Description, "refs/remotes/"+upstream)
		}
	}

	symbol := mod.NoUpstreamSymbol
	aheadBehind := "upToDate"
	if upstream != "" {
		if ahead > 0 && behind > 0 {
			symbol = mod.DivergedSymbol
			aheadBehind = "diverged"
		} else if ahead > 0 {
			symbol = mod.AheadSymbol
			aheadBehind = "ahead"
		} else if behind > 0 {
			symbol = mod.BehindSymbol
			aheadBehind = "behind"
		} else {
			symbol = mod.UpToDateSymbol
			aheadBehind = "upToDate"
		}
	}

	data := gitDivergedResult{
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

func (mod GitDiverged) renderDefault(
	context *Context,
	symbol string,
	data gitDivergedResult,
) string {
	parts := []string{}

	if data.Behind > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", mod.BehindSymbol, data.Behind))
	}
	if data.Ahead > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", mod.AheadSymbol, data.Ahead))
	}
	if data.Behind == 0 && data.Ahead == 0 {
		parts = append(parts, symbol)
	}

	return strings.Join(parts, " ")
}

func init() {
	registerModule(
		"git_diverged",
		registeredModule{
			jsonSchema: schemas.GitDivergedJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := GitDiverged{
					Type:             "git_diverged",
					AheadSymbol:      "↑",
					BehindSymbol:     "↓",
					DivergedSymbol:   "↕",
					UpToDateSymbol:   "≡",
					NoUpstreamSymbol: "?",
				}
				err := node.Decode(&module)
				return &module, err
			},
		},
	)
}
