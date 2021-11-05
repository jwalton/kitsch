package powerline

import (
	"testing"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/kitsch-prompt/internal/styling"
	"github.com/stretchr/testify/assert"
)

func TestPowerline(t *testing.T) {
	gchalk.SetLevel(gchalk.LevelAnsi16m)

	styles := styling.Registry{}
	var powerline = New(&styles, " ", "\ue0b0", " ", false)

	firstSegment := powerline.Segment("red", "hello")
	assert.Equal(t,
		gchalk.BgRed("hello"),
		firstSegment,
	)

	secondSegment := powerline.Segment("blue", "world")

	assert.Equal(t,
		gchalk.WithBgRed().Blue(" ")+gchalk.WithBgBlue().Red("\ue0b0 ")+gchalk.BgBlue("world"),
		secondSegment,
	)

	end := powerline.Finish()

	assert.Equal(t,
		gchalk.WithBgBlue().Black(" ")+gchalk.WithBgBlack().Blue("\ue0b0"),
		end,
	)
}

func TestReversePowerline(t *testing.T) {
	gchalk.SetLevel(gchalk.LevelAnsi16m)

	styles := styling.Registry{}
	var powerline = New(&styles, " ", "\ue0b2", " ", true)

	firstSegment := powerline.Segment("red", "hello")
	assert.Equal(t,
		gchalk.WithBgBlack().Red("\ue0b2")+gchalk.WithBgRed().Black(" ")+gchalk.BgRed("hello"),
		firstSegment,
	)

	secondSegment := powerline.Segment("blue", "world")

	assert.Equal(t,
		gchalk.WithBgRed().Blue(" \ue0b2")+gchalk.WithBgBlue().Red(" ")+gchalk.BgBlue("world"),
		secondSegment,
	)

	end := powerline.Finish()

	assert.Equal(t,
		"",
		end,
	)
}
