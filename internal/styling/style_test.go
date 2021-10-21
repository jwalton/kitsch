package styling

import (
	"testing"

	"github.com/jwalton/gchalk"
	"github.com/stretchr/testify/assert"
)

func testStyleRegistry() Registry {
	gchalkInstance := gchalk.New()
	gchalkInstance.SetLevel(gchalk.LevelAnsi16m)
	styles := Registry{gchalkInstance: gchalkInstance}
	return styles
}

func TestApply(t *testing.T) {
	styles := testStyleRegistry()

	style, err := styles.Get("red")
	assert.NoError(t, err)
	assert.Equal(t, "\u001b[31mtest\u001b[39m", style.Apply("test"))

	style, err = styles.Get("bold")
	assert.NoError(t, err)
	assert.Equal(t, "\u001b[1mtest\u001b[22m", style.Apply("test"))

	style, err = styles.Get("#fff")
	assert.NoError(t, err)
	assert.Equal(t, "\u001b[38;2;255;255;255mtest\u001b[39m", style.Apply("test"))

	style, err = styles.Get("bg:#fff")
	assert.NoError(t, err)
	assert.Equal(t, "\u001b[48;2;255;255;255mtest\u001b[49m", style.Apply("test"))

	styles.AddCustomColor("$foreground", "white")
	style, err = styles.Get("$foreground")
	assert.NoError(t, err)
	assert.Equal(t, "\u001b[37mtest\u001b[39m", style.Apply("test"))
}

// TODO: Add tests for ApplyGetColors

func TestGradient(t *testing.T) {
	styles := testStyleRegistry()

	styles.AddCustomColor("$blue", "#0000ff")
	styles.AddCustomColor("$gradient", "linear-gradient($blue, #fff)")
	_, err := styles.Get("$gradient")
	assert.NoError(t, err)

	styles.AddCustomColor("$gradient2", "linear-gradient($gradient, #fff)")
	_, err = styles.Get("$gradient2")
	assert.EqualError(t,
		err,
		`error compiling style "$gradient2": color $gradient="linear-gradient($blue, #fff)" cannot be used in linear-gradient: invalid hex color "linear-gradient($blue, #fff)"`,
	)

	_, err = styles.Get("linear-gradient(red, blue)")
	assert.EqualError(t,
		err,
		`error compiling style "linear-gradient(red, blue)": expected color at position 1, got "red"`,
	)
}
