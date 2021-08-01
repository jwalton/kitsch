package style

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestParse(t *testing.T) {
	style, err := Parse("blue")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, Style{FG: "blue", BG: "", Modifiers: nil}, style)

	style, err = Parse("bg:blue")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, Style{FG: "", BG: "blue", Modifiers: nil}, style)

	style, err = Parse("bgBlue")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, Style{FG: "", BG: "blue", Modifiers: nil}, style)

	style, err = Parse("red bg:blue bold")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, Style{FG: "red", BG: "blue", Modifiers: []string{"bold"}}, style)

	style, err = Parse("#fff bg:#dead00")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, Style{FG: "#fff", BG: "#dead00", Modifiers: nil}, style)
}

func TestUnmarshall(t *testing.T) {
	var style Style

	err := style.UnmarshalInterface(map[string]interface{}{
		"fg":        "blue",
		"bg":        "red",
		"modifiers": []string{"bold"},
	})
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, Style{FG: "blue", BG: "red", Modifiers: []string{"bold"}}, style)

	err = style.UnmarshalInterface("bgBlue bold")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, Style{FG: "", BG: "blue", Modifiers: []string{"bold"}}, style)
}

type yamlTestStruct struct {
	Style Style `yaml:"style"`
}

func TestUnmarshallYaml(t *testing.T) {
	text := `
style:
  fg: blue
  bg: red
  modifiers:
    - bold
`

	var result yamlTestStruct
	err := yaml.Unmarshal([]byte(text), &result)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, Style{FG: "blue", BG: "red", Modifiers: []string{"bold"}}, result.Style)
}

func TestUnmarshallYamlString(t *testing.T) {
	text := `
style: bgBlue bold
`

	var result yamlTestStruct
	err := yaml.Unmarshal([]byte(text), &result)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, Style{FG: "", BG: "blue", Modifiers: []string{"bold"}}, result.Style)
}
