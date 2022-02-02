package modules

import (
	"fmt"
	"strings"

	"github.com/jwalton/kitsch/internal/gitutils"
	"github.com/jwalton/kitsch/internal/kitsch/modules/schemas"
	"gopkg.in/yaml.v3"
)

//go:generate go run ../genSchema/main.go --pkg schemas GitStateModule

// GitStateModule shows information about the current git repo.
//
type GitStateModule struct {
	// Type is the type of this module.
	Type string `yaml:"type" jsonschema:",required,enum=git_state"`

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

type gitStateResult struct {
	// State is the current state of this repo.
	State gitutils.RepositoryStateType `yaml:"state"`
	// Step is the current step number if we are rebasing, 0 otherwise.
	Step string `yaml:"step"`
	// Total is the total number of steps to complete to finish the rebase, or 0
	// if not rebasing.
	Total string `yaml:"total"`
}

// Execute runs a git module.
func (mod GitStateModule) Execute(context *Context) ModuleResult {
	git := context.Git()

	if git == nil {
		return ModuleResult{DefaultText: "", Data: gitStateResult{}}
	}

	state := git.State()

	data := gitStateResult{
		State: state.State,
		Step:  state.Step,
		Total: state.Total,
	}

	return ModuleResult{
		DefaultText: mod.renderDefault(context, data),
		Data:        data,
	}
}

func (mod GitStateModule) renderDefault(
	context *Context,
	data gitStateResult,
) string {
	if data.State == gitutils.StateNone {
		return ""
	}

	out := strings.Builder{}
	out.WriteString(mod.getStateDescription(data.State))
	if data.Total != "" {
		out.WriteString(fmt.Sprintf(" %s/%s", data.Step, data.Total))
	}

	return out.String()
}

func (mod GitStateModule) getStateDescription(state gitutils.RepositoryStateType) string {
	switch state {
	case gitutils.StateNone:
		return ""
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

func init() {
	registerModule(
		"git_state",
		registeredModule{
			jsonSchema: schemas.GitStateModuleJSONSchema,
			factory: func(node *yaml.Node) (Module, error) {
				module := GitStateModule{
					Type:                "git_state",
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
