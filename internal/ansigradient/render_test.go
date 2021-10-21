package ansigradient

import (
	"testing"

	"github.com/jwalton/gchalk"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	gradient := CSSLinearGradientMust("#ff0020, #2000ff")
	assert.Equal(t,
		"\u001b[38;2;236;0;50mH\u001b[38;2;199;0;87me\u001b[38;2;162;0;124ml\u001b[38;2;124;0;162ml\u001b[38;2;87;0;199mo\u001b[38;2;50;0;236m!\u001b[39m",
		ApplyGradientsRaw("Hello!", gradient, nil, LevelAnsi16m),
	)
	assert.Equal(t,
		"\u001b[48;2;245;0;41mH\u001b[48;2;227;0;59me\u001b[48;2;208;0;78ml\u001b[48;2;189;0;97ml\u001b[48;2;171;0;115mo\u001b[48;2;152;0;134m "+
			"\u001b[48;2;134;0;152mW\u001b[48;2;115;0;171mo\u001b[48;2;97;0;189mr\u001b[48;2;78;0;208ml\u001b[48;2;59;0;227md\u001b[48;2;41;0;245m!\u001b[49m",
		ApplyGradientsRaw("Hello World!", nil, gradient, LevelAnsi16m),
	)
}

func TestRenderPreColoredText(t *testing.T) {
	g := gchalk.New(gchalk.ForceLevel(gchalk.LevelAnsi16m))
	gradient := CSSLinearGradientMust("#ff0000, #ff0000")

	// Render FG with a character that already has a FG color.
	message := "AB" + g.Green("C") + "DE"
	output := ApplyGradientsRaw(message, gradient, nil, LevelAnsi16m)
	assert.Equal(t, "\u001b[38;2;255;0;0mAB\u001b[32mC\u001b[38;2;255;0;0mDE\u001b[39m", output)

	// Render BG with a character that has a FG color.
	output = ApplyGradientsRaw(message, nil, gradient, LevelAnsi16m)
	assert.Equal(t, "\u001b[48;2;255;0;0mAB\u001b[32mC\u001b[39mDE\u001b[49m", output)

	// Render BG with a character that already has a BG color.
	message = "AB" + g.BgGreen("C") + "DE"
	output = ApplyGradientsRaw(message, nil, gradient, LevelAnsi16m)
	assert.Equal(t, "\u001b[48;2;255;0;0mAB\u001b[42mC\u001b[48;2;255;0;0mDE\u001b[49m", output)

	// Render FG with a character that has a BG color.
	output = ApplyGradientsRaw(message, gradient, nil, LevelAnsi16m)
	assert.Equal(t, "\u001b[38;2;255;0;0mAB\u001b[42mC\u001b[49mDE\u001b[39m", output)
}

func TestRenderGradientsOverGradients(t *testing.T) {
	bgGradient := CSSLinearGradientMust("#f00, #000")
	fgGradient := CSSLinearGradientMust("#000, #0ff")

	fgBgExpected := "\u001b[38;2;0;10;10m\u001b[48;2;244;0;0mH\u001b[38;2;0;31;31m\u001b[48;2;223;0;0me\u001b[38;2;0;53;53m\u001b[48;2;201;0;0ml\u001b[38;2;0;74;74m\u001b[48;2;180;0;0ml\u001b[38;2;0;95;95m\u001b[48;2;159;0;0mo\u001b[38;2;0;116;116m\u001b[48;2;138;0;0m " +
		"\u001b[38;2;0;138;138m\u001b[48;2;116;0;0mw\u001b[38;2;0;159;159m\u001b[48;2;95;0;0mo\u001b[38;2;0;180;180m\u001b[48;2;74;0;0mr\u001b[38;2;0;201;201m\u001b[48;2;53;0;0ml\u001b[38;2;0;223;223m\u001b[48;2;31;0;0md\u001b[38;2;0;244;244m\u001b[48;2;10;0;0m!\u001b[39m\u001b[49m"
	bgFgExpected := "\u001b[48;2;244;0;0m\u001b[38;2;0;10;10mH\u001b[48;2;223;0;0m\u001b[38;2;0;31;31me\u001b[48;2;201;0;0m\u001b[38;2;0;53;53ml\u001b[48;2;180;0;0m\u001b[38;2;0;74;74ml\u001b[48;2;159;0;0m\u001b[38;2;0;95;95mo\u001b[48;2;138;0;0m\u001b[38;2;0;116;116m " +
		"\u001b[48;2;116;0;0m\u001b[38;2;0;138;138mw\u001b[48;2;95;0;0m\u001b[38;2;0;159;159mo\u001b[48;2;74;0;0m\u001b[38;2;0;180;180mr\u001b[48;2;53;0;0m\u001b[38;2;0;201;201ml\u001b[48;2;31;0;0m\u001b[38;2;0;223;223md\u001b[48;2;10;0;0m\u001b[38;2;0;244;244m!\u001b[49m\u001b[39m"

	output := ApplyGradientsRaw("Hello world!", fgGradient, bgGradient, LevelAnsi16m)
	assert.Equal(t, fgBgExpected, output)

	fgOnly := ApplyGradientsRaw("Hello world!", fgGradient, nil, LevelAnsi16m)
	fgFirstOutput := ApplyGradientsRaw(fgOnly, nil, bgGradient, LevelAnsi16m)
	assert.Equal(t, fgBgExpected, fgFirstOutput)

	bgOnly := ApplyGradientsRaw("Hello world!", nil, bgGradient, LevelAnsi16m)
	bgFirstOutput := ApplyGradientsRaw(bgOnly, fgGradient, nil, LevelAnsi16m)
	assert.Equal(t, bgFgExpected, bgFirstOutput)
}

func TestRenderMultipleCharsSameColor(t *testing.T) {
	gradient := CSSLinearGradientMust("#ff0000 3px, #0000ff 3px")
	assert.Equal(t,
		"\u001b[48;2;255;0;0mRed\u001b[48;2;0;0;255mBlu\u001b[49m",
		ApplyGradientsRaw("RedBlu", nil, gradient, LevelAnsi16m),
	)
}

func TestZwjEmojiRender(t *testing.T) {
	gradient := CSSLinearGradientMust("#ff0000, #0000ff")

	// This string has an astronaut made up of woman, light skin tone, rocket,
	// and the zero-width-joiner codepoints, all combined into a single
	// grapheme.  If we use a naive "insert escape codes between each rune"
	// approach, we'd end up splitting this into its component emojis.
	str := "AB👩🏻‍🚀"

	// FIXME: This is wrong - the astronaut should be 2 characters wide.
	assert.Equal(t,
		"\x1b[48;2;212;0;42mA\x1b[48;2;127;0;127mB\x1b[48;2;42;0;212m👩🏻\u200d🚀\x1b[49m",
		ApplyGradientsRaw(str, nil, gradient, LevelAnsi16m),
	)
}

func Test256ColorMode(t *testing.T) {
	gradient := CSSLinearGradientMust("#ff0020 0%, #2000ff 100%")

	assert.Equal(t,
		"\u001b[38;5;197mH\u001b[38;5;161me\u001b[38;5;162mll\u001b[38;5;126mo\u001b[38;5;127m W\u001b[38;5;91mo\u001b[38;5;92mrl\u001b[38;5;56md\u001b[38;5;57m!\u001b[39m",
		ApplyGradientsRaw("Hello World!", gradient, nil, LevelAnsi256),
	)

	assert.Equal(t,
		"\u001b[48;5;197mH\u001b[48;5;161me\u001b[48;5;162mll\u001b[48;5;126mo\u001b[48;5;127m W\u001b[48;5;91mo\u001b[48;5;92mrl\u001b[48;5;56md\u001b[48;5;57m!\u001b[49m",
		ApplyGradientsRaw("Hello World!", nil, gradient, LevelAnsi256),
	)
}
