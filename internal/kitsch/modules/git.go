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
// Provides the following template variables:
//
// • State - A `{ State, Step, Total, Base, Branch }` object.  All of these
//   values are strings.  State is the current state of the repo (e.g. "MERGING"
//   if in the middle of a merge).  Step and Total represent the number of steps
//   left to complete the current operation (e.g. the number of commits left
//   to apply in an interactive rebase), or empty string if no such operation is
//   in progress.  Base is the name of the base branch we are merging from or
//   rebasing from.  Branch is the name of the current branch or hash.
//
// • Ahead - The number of commits ahead of the upstream branch.
//
// • Behind - The number of commits behind the upstream branch.
//
// • Symbol - The symbol to use to indicate the current state of the repo.
//
type GitModule struct {
	CommonConfig `yaml:",inline"`
	// Type is the type of this module.
	Type             string `yaml:"type" jsonschema:",enum=git"`
	AheadStyle       string `yaml:"aheadStyle"`
	BehindStyle      string `yaml:"behindStyle"`
	AheadBehindStyle string `yaml:"aheadBehindStyle"`
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

	style := defaultString(mod.Style, "brightCyan")
	if upstream != "" {
		if ahead > 0 && behind > 0 {
			style = defaultString(mod.AheadBehindStyle, "brightYellow")
		} else if ahead > 0 {
			style = defaultString(mod.AheadStyle, "brightGreen")
		} else if behind > 0 {
			style = defaultString(mod.BehindStyle, "brightRed")
		}
	}

	data := map[string]interface{}{
		"State":  state,
		"Ahead":  ahead,
		"Behind": behind,
		"Symbol": symbol,
	}

	defaultOutput := mod.renderDefault(context, symbol, state, upstream, ahead, behind)

	return executeModule(context, mod.CommonConfig, data, style, defaultOutput)
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
