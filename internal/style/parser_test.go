package style

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStyle(t *testing.T) {
	customColors := map[string]string{}
	style, err := parseStyle(customColors, "blue")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "blue", bg: "", modifiers: nil}, style)

	style, err = parseStyle(customColors, "bg:blue")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "", bg: "blue", modifiers: nil}, style)

	style, err = parseStyle(customColors, "bgBlue")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "", bg: "blue", modifiers: nil}, style)

	style, err = parseStyle(customColors, "red bg:blue bold")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "red", bg: "blue", modifiers: []string{"bold"}}, style)

	style, err = parseStyle(customColors, "#fff bg:#dead00")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "#fff", bg: "#dead00", modifiers: nil}, style)

	style, err = parseStyle(customColors, "linear-gradient(#fff, #000) bg:linear-gradient(#f00, #00F)")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "linear-gradient(#fff, #000)", bg: "linear-gradient(#f00, #00F)", modifiers: nil}, style)

	style, err = parseStyle(customColors, "red green grey")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "grey", bg: "", modifiers: nil}, style)

	_, err = parseStyle(customColors, "banana")
	assert.EqualError(t, err, "unknown style \"banana\"")

	_, err = parseStyle(customColors, "bg:banana")
	assert.EqualError(t, err, "unknown style \"bg:banana\"")
}

func TestCustomColors(t *testing.T) {
	customColors := map[string]string{
		"$foreground": "#fff",
		"$background": "#000",
		"$gradient":   "linear-gradient(#f00, #00f)",
		"$red":        "blue",
	}

	style, err := parseStyle(customColors, "$red")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "blue", bg: "", modifiers: nil}, style)

	style, err = parseStyle(customColors, "bg:$red")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "", bg: "blue", modifiers: nil}, style)

	style, err = parseStyle(customColors, "$foreground bg:$background")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "#fff", bg: "#000", modifiers: nil}, style)

	style, err = parseStyle(customColors, "$gradient")
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, styleDescriptor{fg: "linear-gradient(#f00, #00f)", bg: "", modifiers: nil}, style)

	_, err = parseStyle(customColors, "$banana")
	assert.EqualError(t, err, "unknown style \"$banana\"")
}
