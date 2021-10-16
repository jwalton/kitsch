package style

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
