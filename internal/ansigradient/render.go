package ansigradient

import (
	"image/color"
	"strings"

	"github.com/jwalton/gchalk/pkg/ansistyles"
)

// ApplyGradients will apply the given background and foreground gradients to the given string.
//
// ApplyGradients will attempt to automatically detect the current color support level
// based on stdout - this can be overridden globally by calling `SetLevel()` or `ApplyGradientsRaw()`.
func ApplyGradients(str string, foreground Gradient, background Gradient) string {
	result, _ := ApplyGradientsRawLen(str, foreground, background, GetLevel())
	return result
}

// ApplyGradientsRaw will apply the given background and foreground gradients to the given string.
// Unlike ApplyGradients, this will not attempt to automatically detect the current color support level.
// Returns the colored string, and the print width of the string.
func ApplyGradientsRaw(str string, foreground Gradient, background Gradient, level ColorLevel) string {
	result, _ := ApplyGradientsRawLen(str, foreground, background, level)
	return result
}

// ApplyGradientsLen will apply the given background and foreground gradients to the given string,
// and return the "print width" of the string in columns.
//
// ApplyGradients will attempt to automatically detect the current color support level
// based on stdout - this can be overridden globally by calling `SetLevel()` or `ApplyGradientsRaw()`.
func ApplyGradientsLen(str string, foreground Gradient, background Gradient) (string, int) {
	return ApplyGradientsRawLen(str, foreground, background, GetLevel())
}

// ApplyGradientsRawLen will apply the given background and foreground gradients to the given string,
// and return the "print width" of the string in columns.
//
// Unlike ApplyGradients, this will not attempt to automatically detect the current color support level.
func ApplyGradientsRawLen(str string, foreground Gradient, background Gradient, level ColorLevel) (string, int) {
	parsed, printWidth := tokenize(str)

	if foreground == nil && background == nil {
		return str, printWidth
	}

	var fgColors ColorGenerator
	var bgColors ColorGenerator
	if foreground != nil {
		// TODO: Rename this
		fgColors = foreground.Generator(printWidth)
	}
	if background != nil {
		bgColors = background.Generator(printWidth)
	}

	// TODO: Rename this
	return renderRGBAs(parsed, printWidth, fgColors, bgColors, level), printWidth
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
	parsed []gradientToken,
	printWidth int,
	fgColors ColorGenerator,
	bgColors ColorGenerator,
	level ColorLevel,
) string {
	out := strings.Builder{}
	worstCaseLength := printWidth * 20
	out.Grow(worstCaseLength)

	column := 0
	context := colorizeContext{level: level}

	prevFg := ""
	prevBg := ""

	for _, token := range parsed {
		switch token.t {
		case tokenString:
			if (fgColors == nil || token.fg != "") && (bgColors == nil || token.bg != "") {
				// Don't color this string.
				out.Write([]byte(token.content))
			} else {
				colorizeASCIIString(&context, token.content, column, fgColors, bgColors, &out)
			}
			column += len(token.content)

		case tokenComplexChar:
			if (bgColors == nil || token.fg != "") && (bgColors == nil || token.bg != "") {
				// Don't color this string.
			} else {
				renderColorCodes(&context, float64(column)+(float64(token.printWidth)/2), fgColors, bgColors, &out)
			}
			out.WriteString(token.content)

			column += token.printWidth

		case tokenEscapeCode:
			if fgColors != nil && token.content == ansistyles.Close {
				// Skip closing the forground color.
				context.lastFgColor = color.RGBA{}
			} else if bgColors != nil && token.content == ansistyles.BgClose {
				// Skip closing the background color.
				context.lastBgColor = color.RGBA{}
			} else {
				// Reset the last color codes.
				if prevFg != token.fg {
					context.lastFgColor = color.RGBA{}
				}
				if prevBg != token.bg {
					context.lastBgColor = color.RGBA{}
				}

				// Pass the escape code through
				out.WriteString(token.content)
			}

			prevFg = token.fg
			prevBg = token.bg

		default:
			// Pass this through
			out.WriteString(token.content)
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
	column int,
	fgColors ColorGenerator,
	bgColors ColorGenerator,
	out *strings.Builder,
) {
	for charIndex := range str {
		renderColorCodes(context, float64(column)+0.5, fgColors, bgColors, out)
		out.WriteByte(str[charIndex])
		column++
	}
}

// Writes escape sequences for the foreground and background colors to `out`.
func renderColorCodes(
	context *colorizeContext,
	position float64,
	fgColors ColorGenerator,
	bgColors ColorGenerator,
	out *strings.Builder,
) {
	// Write the forground color, if any.
	if fgColors != nil {
		nextColor := fgColors.ColorAt(position)
		if context.lastFgColor.A == 0 || !compareColors[context.level](nextColor, context.lastFgColor) {
			context.lastFgColor = nextColor
			fgColorize[context.level](out, nextColor)
		}
	}

	// Write the background color, if any.
	if bgColors != nil {
		nextColor := bgColors.ColorAt(position)
		if context.lastBgColor.A == 0 || !compareColors[context.level](nextColor, context.lastBgColor) {
			context.lastBgColor = nextColor
			bgColorize[context.level](out, nextColor)
		}
	}
}
