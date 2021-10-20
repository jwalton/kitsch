// Package style is for parsing style strings from the kitsch configuration,
// and for applying styles to strings.
package styling

import (
	"strings"

	"github.com/jwalton/gchalk"
	"github.com/jwalton/go-ansiparser"
	"github.com/jwalton/kitsch-prompt/internal/ansigradient"
	"github.com/jwalton/kitsch-prompt/internal/colortools"
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

		if strings.HasPrefix(token, linearGradientPrefix) {
			cssGradient := token[len(linearGradientPrefix) : len(token)-1]
			if background {
				bgGradient, err = ansigradient.CSSLinearGradientWithMap(customColors, cssGradient)
			} else {
				fgGradient, err = ansigradient.CSSLinearGradientWithMap(customColors, cssGradient)
			}
		} else {
			if background {
				// TODO: Must be a better way!  :/
				token = "bg" + strings.ToUpper(token[0:1]) + token[1:]
			}
			builder, err = builder.WithStyle(token)
		}

		return err
	}

	err = compileColor(descriptor.fg, false)
	if err != nil {
		return Style{}, err
	}

	err = compileColor(descriptor.bg, true)
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
	if style == nil {
		return text
	}
	if style.builder != nil {
		text = style.builder.Paint(text)
	}
	text = ansigradient.ApplyGradientsRaw(text, style.fgGradient, style.bgGradient, gchalk.GetLevel())
	return text
}

// ApplyGetColors applies this style to the given text, and returns the first and last colors of the styled text.
func (style *Style) ApplyGetColors(text string) (result string, first CharacterColors, last CharacterColors) {
	if style == nil {
		return text, first, last
	}
	result = style.Apply(text)

	// TODO: This is not very efficient, because we end up parsing the string
	// twice.  Consider making it so `ApplyGradients` lets us pass in a
	// pre-parsed string?
	printLength := 0
	if style.fgGradient != nil || style.bgGradient != nil {
		parsed := ansiparser.Parse(text)
		printLength = ansiparser.TokensPrintLength(parsed)
	}

	first.FG, last.FG = getCharacterColors(style.descriptor.fg, style.fgGradient, printLength)
	first.BG, last.BG = getCharacterColors(style.descriptor.bg, style.bgGradient, printLength)

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
