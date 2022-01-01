package config

import (
	"testing"

	"github.com/jwalton/kitsch-prompt/sampleconfig"
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
