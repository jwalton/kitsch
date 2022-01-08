// Package styling is for parsing style strings from the kitsch configuration,
// and for applying styles to strings.
package styling

import (
	"image/color"
	"strings"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/gchalk/pkg/ansistyles"
	"github.com/jwalton/kitsch/internal/ansigradient"
	"github.com/jwalton/kitsch/internal/colortools"
)

// Style represents a compiled style that can be applied to a string.
// Styles are constructed by calling into StyleRegistry.Get().  Once created,
// a Style is immutable.
type Style struct {
	descriptor styleDescriptor
	builder    *gchalk.Builder
	fgGradient ansigradient.Gradient
	bgGradient ansigradient.Gradient
}

// CharacterColors represent the color for a single character.
type CharacterColors struct {
	// FG is the foreground color.  This is either a color name like "red", or
	// a hex color.
	FG string
	// BG is the background color.
	BG string
}

func compileStyle(
	baseBuilder *gchalk.Builder,
	customColors map[string]string,
	styleString string,
) (Style, error) {
	descriptor, err := parseStyle(customColors, styleString)
	if err != nil {
		return Style{}, err
	}

	builder := baseBuilder

	var fgGradient ansigradient.Gradient
	var bgGradient ansigradient.Gradient

	compileColor := func(token string, background bool) error {
		var err error

		if token == "" {
			return nil
		}

		if _, validAnsiStyle := ansistyles.Color[token]; validAnsiStyle {
			if background {
				// TODO: Must be a better way!  :/
				token = "bg" + strings.ToUpper(token[0:1]) + token[1:]
			}
			builder, err = builder.WithStyle(token)
		} else if strings.HasPrefix(token, linearGradientPrefix) {
			cssGradient := token[len(linearGradientPrefix) : len(token)-1]
			if background {
				bgGradient, err = ansigradient.CSSLinearGradientWithMap(customColors, cssGradient)
			} else {
				fgGradient, err = ansigradient.CSSLinearGradientWithMap(customColors, cssGradient)
			}
		} else {
			var c color.RGBA
			c, err = colortools.ParseColor(token)
			if background {
				builder = builder.WithBgRGB(c.R, c.G, c.B)
			} else {
				builder = builder.WithRGB(c.R, c.G, c.B)
			}
		}

		return err
	}

	err = compileColor(descriptor.bg, true)
	if err != nil {
		return Style{}, err
	}

	err = compileColor(descriptor.fg, false)
	if err != nil {
		return Style{}, err
	}

	for _, modifier := range descriptor.modifiers {
		builder, err = builder.WithStyle(modifier)
		if err != nil {
			return Style{}, err
		}
	}

	return Style{
		descriptor: descriptor,
		builder:    builder,
		fgGradient: fgGradient,
		bgGradient: bgGradient,
	}, nil
}

// Apply applies this style to the given text.
func (style *Style) Apply(text string) string {
	result, _, _ := style.ApplyGetColors(text)
	return result
}

// ApplyGetColors applies this style to the given text, and returns the first and last colors of the styled text.
func (style *Style) ApplyGetColors(text string) (result string, first CharacterColors, last CharacterColors) {
	if style == nil {
		return text, first, last
	}

	if style.builder != nil {
		result = style.builder.Paint(text)
	}

	printWidth := 0
	if style.fgGradient != nil || style.bgGradient != nil {
		// TODO: This instance of gchalk is not the same instance as the one from the styleRegistry.
		result, printWidth = ansigradient.ApplyGradientsRawLen(text, style.fgGradient, style.bgGradient, gchalk.GetLevel())
	}

	first.FG, last.FG = getCharacterColors(style.descriptor.fg, style.fgGradient, printWidth)
	first.BG, last.BG = getCharacterColors(style.descriptor.bg, style.bgGradient, printWidth)
	if first.BG != "" {
		first.BG = "bg:" + first.BG
	}
	if last.BG != "" {
		last.BG = "bg:" + last.BG
	}

	return result, first, last
}

func getCharacterColors(colorString string, gradient ansigradient.Gradient, printLength int) (first string, last string) {
	if gradient != nil {
		first = colortools.ColorToHex(gradient.ColorAt(printLength, -1))
		last = colortools.ColorToHex(gradient.ColorAt(printLength, printLength+1))
	} else {
		first = colorString
		last = colorString
	}
	return first, last
}
