// Package style is for parsing style strings from the kitsch configuration,
// and for applying styles to strings.
package style

import (
	"strings"

	"github.com/jwalton/gchalk"
)

type stylerFn func(string) string

// Style represents a compiled style that can be applied to a string.
// Styles are constructed by calling into StyleRegistry.Get().  Once created,
// a Style is immutable.
type Style struct {
	descriptor styleDescriptor
	styler     stylerFn
	builder    *gchalk.Builder
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
	var styler stylerFn
	builder := baseBuilder

	compileColor := func(token string, background bool) error {
		var err error

		if token == "" {
			return nil
		}

		if strings.HasPrefix(token, linearGradientPrefix) {
			// TODO - linear-gradient
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
		styler:     styler,
		builder:    builder,
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
	if style.styler != nil {
		text = style.styler(text)
	}
	return text
}

// ApplyGetColors applies this style to the given text, and returns the first and last colors of the styled text.
func (style *Style) ApplyGetColors(text string) (result string, first CharacterColors, last CharacterColors) {
	if style == nil {
		return text, first, last
	}
	if style.builder != nil {
		text = style.builder.Paint(text)
	}
	if style.styler != nil {
		text = style.styler(text)
	}

	if strings.HasPrefix(style.descriptor.fg, linearGradientPrefix) {
		// TODO - linear-gradient
	} else {
		first.FG = style.descriptor.fg
		last.FG = style.descriptor.fg
	}

	if strings.HasPrefix(style.descriptor.bg, linearGradientPrefix) {
		// TODO - linear-gradient
	} else {
		first.BG = style.descriptor.bg
		last.BG = style.descriptor.bg
	}

	return text, first, last
}
