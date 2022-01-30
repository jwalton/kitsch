package modules

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	blockMod := moduleWrapperFromYAML(heredoc.Doc(`
		type: block
		modules:
		- type: text
		  text: hello
		- type: text
		  text: world
    `))

	result := blockMod.Execute(newTestContext("jwalton"))
	assert.Equal(t, "hello world", result.Text)
}

func TestBlockStyles(t *testing.T) {
	blockMod := moduleWrapperFromYAML(heredoc.Doc(`
		type: block
		join: " {{.PrevColors.FG}}{{.NextColors.FG}} "
		modules:
		- type: text
		  style: red
		  text: hello
		- type: text
		  style: blue
		  text: world
    `))

	result := blockMod.Execute(newTestContext("jwalton"))
	assert.Equal(t, "hello redblue world", result.Text)
}

// TestBlockSubIDs verifies that the results of child modules can be indexed by ID.
func TestBlockSubIDs(t *testing.T) {
	blockMod := moduleWrapperFromYAML(heredoc.Doc(`
		type: block
		template: "{{ with .Data.Modules.a }}{{ .Data.Username }}{{ end }}"
		modules:
		- type: username
		  id: a
		  showAlways: true
		- type: prompt
    `))

	result := blockMod.Execute(newTestContext("oriana"))

	assert.Equal(t, "oriana", result.Text)
}
