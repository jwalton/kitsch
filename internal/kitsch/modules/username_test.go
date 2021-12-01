package modules

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/env"
	"github.com/stretchr/testify/assert"
)

func TestUsernameNoSSH(t *testing.T) {
	mod := moduleFromYAMLMust(heredoc.Doc(`
		type: username
	`))
	context := newTestContext("jwalton")

	result := mod.Execute(context)
	assert.Equal(t, "", result.Text)
}

func TestUsernameNoSSHWithTemplate(t *testing.T) {
	mod := moduleFromYAMLMust(heredoc.Doc(`
		type: username
		template: '{{ .Data.Username }}'
	`))
	context := newTestContext("jwalton")

	result := mod.Execute(context)
	assert.Equal(t, "jwalton", result.Text)
}

func TestUsername(t *testing.T) {
	mod := moduleFromYAMLMust(heredoc.Doc(`
		type: username
	`))

	context := newTestContext("jwalton")
	context.Environment = &env.DummyEnv{
		Env: map[string]string{
			"USER":    "jwalton",
			"HOME":    "/Users/jwalton",
			"SSH_TTY": "true",
		},
	}

	result := mod.Execute(context)

	assert.Equal(t, "jwalton", result.Text)
	assert.Equal(t,
		usernameModuleData{
			username: "jwalton",
			IsSSH:    true,
			Show:     true,
		},
		result.Data.(usernameModuleData),
	)
	assert.Equal(t, "jwalton", result.Data.(usernameModuleData).Username())
}
