package ansigradient

import (
	"image/color"
	"strings"

	"github.com/jwalton/gchalk/pkg/ansistyles"
	"github.com/jwalton/go-ansiparser"
)

// ApplyGradients will apply the given background and foreground gradients to the given string.
//
// ApplyGradients will attempt to automatically detect the current color support level
// based on stdout - this can be overridden globally by calling `SetLevel()` or `ApplyGradientsRaw()`.
func ApplyGradients(str string, foreground Gradient, background Gradient) string {
	return ApplyGradientsRaw(str, foreground, background, GetLevel())
}

// ApplyGradientsRaw will apply the given background and foreground gradients to the given string.
// Unlike ApplyGradients, this will not attempt to automatically detect the current color support level.
func ApplyGradientsRaw(str string, foreground Gradient, background Gradient, level ColorLevel) string {
	if foreground == nil && background == nil {
		return str
	}

	parsed := ansiparser.Parse(str)
	printLength := ansiparser.TokensPrintLength(parsed)

	var fgColors []color.RGBA
	var bgColors []color.RGBA
	if foreground != nil {
		fgColors = foreground.Colors(printLength)
	}
	if background != nil {
		bgColors = background.Colors(printLength)
	}

	return renderRGBAs(parsed, fgColors, bgColors, level)
}

// ColorString will apply the given background and foreground colors to the given string.
//
// ColorString will attempt to automatically detect the current color support level
// based on stdout - this can be overridden globally by calling `SetLevel()` or `ColorizeStringRaw()`.
func ColorString(str string, foreground []color.RGBA, background []color.RGBA) string {
	return ColorStringRaw(str, foreground, background, GetLevel())
}

// ColorStringRaw will apply the given background and foreground colors to the given string.
// Unlike ColorString, this will not attempt to automatically detect the current color support level.
func ColorStringRaw(str string, foreground []color.RGBA, background []color.RGBA, level ColorLevel) string {
	parsed := ansiparser.Parse(str)
	return renderRGBAs(parsed, foreground, background, level)
}

var compareColors = [](func(color.RGBA, color.RGBA) bool){
	LevelNone: func(a, b color.RGBA) bool { return true },
	LevelBasic: func(a, b color.RGBA) bool {
		return ansistyles.Ansi256ToAnsi(ansistyles.RGBToAnsi256(a.R, a.G, a.B)) ==
			ansistyles.Ansi256ToAnsi(ansistyles.RGBToAnsi256(b.R, b.G, b.B))
	},
	LevelAnsi256: func(a, b color.RGBA) bool {
		return ansistyles.RGBToAnsi256(a.R, a.G, a.B) == ansistyles.RGBToAnsi256(b.R, b.G, b.B)
	},
	LevelAnsi16m: func(a, b color.RGBA) bool {
		return a.R == b.R && a.G == b.G && a.B == b.B
	},
}

var fgColorize = [](func(out *strings.Builder, c color.RGBA)){
	LevelNone: func(out *strings.Builder, c color.RGBA) {},
	LevelBasic: func(out *strings.Builder, c color.RGBA) {
		ansistyles.WriteStringAnsi(out, ansistyles.Ansi256ToAnsi(ansistyles.RGBToAnsi256(c.R, c.G, c.B)))
	},
	LevelAnsi256: func(out *strings.Builder, c color.RGBA) {
		ansistyles.WriteStringAnsi256(out, ansistyles.RGBToAnsi256(c.R, c.G, c.B))
	},
	LevelAnsi16m: func(out *strings.Builder, c color.RGBA) {
		ansistyles.WriteStringAnsi16m(out, c.R, c.G, c.B)
	},
}

var bgColorize = [](func(out *strings.Builder, c color.RGBA)){
	LevelNone: func(out *strings.Builder, c color.RGBA) {},
	LevelBasic: func(out *strings.Builder, c color.RGBA) {
		ansistyles.WriteStringBgAnsi(out, ansistyles.Ansi256ToAnsi(ansistyles.RGBToAnsi256(c.R, c.G, c.B)))
	},
	LevelAnsi256: func(out *strings.Builder, c color.RGBA) {
		ansistyles.WriteStringBgAnsi256(out, ansistyles.RGBToAnsi256(c.R, c.G, c.B))
	},
	LevelAnsi16m: func(out *strings.Builder, c color.RGBA) {
		ansistyles.WriteStringBgAnsi16m(out, c.R, c.G, c.B)
	},
}

type colorizeContext struct {
	level       ColorLevel
	lastFgColor color.RGBA
	lastBgColor color.RGBA
}

// renderRGBAs will take a parsed input string and colors, and write the colorized
// version to the `out` writer.
func renderRGBAs(
	parsed []ansiparser.AnsiToken,
	fgColors []color.RGBA,
	bgColors []color.RGBA,
	level ColorLevel,
) string {
	out := strings.Builder{}
	worstCaseLength := ansiparser.TokensPrintLength(parsed) * 20
	out.Grow(worstCaseLength)

	// Which color to use for the next character.
	colorIndex := 0
	context := colorizeContext{level: level}

	if len(fgColors) == 0 {
		fgColors = nil
	}
	if len(bgColors) == 0 {
		bgColors = nil
	}

	prevFg := ""
	prevBg := ""

	for _, token := range parsed {
		switch token.Type {
		case ansiparser.String:
			if (fgColors == nil || token.FG != "") && (bgColors == nil || token.BG != "") {
				// Don't color this string.
				out.Write([]byte(token.Content))
			} else {
				colorizeASCIIString(&context, token.Content, colorIndex, fgColors, bgColors, &out)
			}
			colorIndex += len(token.Content)

		case ansiparser.ComplexChar:
			if (bgColors == nil || token.FG != "") && (bgColors == nil || token.BG != "") {
				// Don't color this string.
			} else {
				renderColorCodes(&context, colorIndex, fgColors, bgColors, &out)
			}
			out.WriteString(token.Content)

			colorIndex++

		case ansiparser.EscapeCode:
			if fgColors != nil && token.Content == ansistyles.Close {
				// Skip closing the forground color.
				context.lastFgColor = color.RGBA{}
			} else if bgColors != nil && token.Content == ansistyles.BgClose {
				// Skip closing the background color.
				context.lastBgColor = color.RGBA{}
			} else {
				// Reset the last color codes.
				if prevFg != token.FG {
					context.lastFgColor = color.RGBA{}
				}
				if prevBg != token.BG {
					context.lastBgColor = color.RGBA{}
				}

				// Pass the escape code through
				out.WriteString(token.Content)
			}

			prevFg = token.FG
			prevBg = token.BG

		default:
			// Pass this through
			out.WriteString(token.Content)
		}
	}

	if fgColors != nil {
		out.WriteString(ansistyles.Close)
	}

	if bgColors != nil {
		out.WriteString(ansistyles.BgClose)
	}

	return out.String()
}

// colorizeASCIIString will colorize a string consisting of ASCII characters,
// where each character is a single byte long.
//
// fgColors is the (possibly 0-length) slice of colors to apply, and bgColors
// is the same but for the background.
func colorizeASCIIString(
	context *colorizeContext,
	str string,
	colorIndex int,
	fgColors []color.RGBA,
	bgColors []color.RGBA,
	out *strings.Builder,
) {
	for charIndex := range str {
		renderColorCodes(context, colorIndex, fgColors, bgColors, out)
		out.WriteByte(str[charIndex])
		colorIndex++
	}
}

// Writes escape sequences for the foreground and background colors to `out`.
func renderColorCodes(
	context *colorizeContext,
	colorIndex int,
	fgColors []color.RGBA,
	bgColors []color.RGBA,
	out *strings.Builder,
) {
	// Write the forground color, if any.
	if colorIndex < len(fgColors) {
		nextColor := fgColors[colorIndex]
		if context.lastFgColor.A == 0 || !compareColors[context.level](nextColor, context.lastFgColor) {
			context.lastFgColor = nextColor
			fgColorize[context.level](out, nextColor)
		}
	}

	// Write the background color, if any.
	if colorIndex < len(bgColors) {
		nextColor := bgColors[colorIndex]
		if context.lastBgColor.A == 0 || !compareColors[context.level](nextColor, context.lastBgColor) {
			context.lastBgColor = nextColor
			bgColorize[context.level](out, nextColor)
		}
	}
}
