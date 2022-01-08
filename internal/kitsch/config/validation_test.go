package config

import (
	"fmt"
	"testing"

	"github.com/jwalton/kitsch/sampleconfig"
	"github.com/stretchr/testify/assert"
)

func TestValidateSimpleConfig(t *testing.T) {
	c := `
prompt:
  type: text
  text: "Hello, world!"
  style: blue
`
	err := ValidateConfiguration([]byte(c))
	assert.Nil(t, err)
}

func TestValidateSimpleConfigWithError(t *testing.T) {
	c := `prompt:
  type: text
  text: "Hello, world!"
  style: blue
  foo: bar
`
	err := ValidateConfiguration([]byte(c))
	fmt.Println(err.Error())
	// TODO: Make this error message better.
	assert.Contains(t, err.Error(), "does not validate")
	// assert.Contains(t, err.Error(), "text (2:3)")
	// assert.Contains(t, err.Error(), "additionalProperties 'foo' not allowed")
}

func TestValidateConfigWithBlock(t *testing.T) {
	c := `
prompt:
  type: block
  modules:
    - type: project
      style: brightBlack
`
	err := ValidateConfiguration([]byte(c))
	assert.Nil(t, err)
}

func TestValidateBuiltInConfigs(t *testing.T) {
	err := ValidateConfiguration(sampleconfig.DefaultConfig)
	assert.Nil(t, err)
}
