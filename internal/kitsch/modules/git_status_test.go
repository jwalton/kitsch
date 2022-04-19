package modules

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch/internal/gitutils"
	"github.com/jwalton/kitsch/internal/kitsch/styling"
	"github.com/stretchr/testify/assert"
)

func TestGitStatusClean(t *testing.T) {
	context := NewDemoContext(
		DemoConfig{},
		&styling.Registry{},
	)

	mod := moduleWrapperFromYAML(heredoc.Doc(`
		type: git_status
	`))

	result := mod.Execute(&context)
	assert.Equal(t, "", result.Text)
}

func TestGitStatusChanges(t *testing.T) {
	context := NewDemoContext(
		DemoConfig{
			Git: gitutils.DemoGit{
				CurrentStats: gitutils.GitStats{
					Index: gitutils.GitFileStats{
						Added:    1,
						Modified: 2,
						Deleted:  3,
					},
					Unmerged: 4,
					Unstaged: gitutils.GitFileStats{
						Added:    5,
						Modified: 6,
						Deleted:  7,
					},
				},
				StashCount: 8,
			},
		},
		&styling.Registry{},
	)

	mod := moduleWrapperFromYAML(heredoc.Doc(`
		type: git_status
	`))

	result := mod.Execute(&context)
	assert.Equal(t, "+1 ~2 -3 !4 | +5 ~6 -7 (8)", result.Text)
}

func TestGitStatusUnmerged(t *testing.T) {
	context := NewDemoContext(
		DemoConfig{
			Git: gitutils.DemoGit{
				CurrentStats: gitutils.GitStats{
					Unmerged: 4,
				},
			},
		},
		&styling.Registry{},
	)

	mod := moduleWrapperFromYAML(heredoc.Doc(`
		type: git_status
	`))

	result := mod.Execute(&context)
	assert.Equal(t, "+0 ~0 -0 !4", result.Text)
}
