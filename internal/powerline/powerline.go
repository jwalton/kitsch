package powerline

import (
	"text/template"

	"github.com/jwalton/kitsch-prompt/internal/styling"
)

// Powerline is a helper object for constructing powerline prompts.
type Powerline struct {
	styles          *styling.Registry
	lastColor       styling.CharacterColors
	separatorPrefix string
	separator       string
	separatorSuffix string
}

// New creates a noew Powerline helper object for use in a template.
func New(styles *styling.Registry, prefix string, separator string, suffix string) *Powerline {
	return &Powerline{
		styles:          styles,
		lastColor:       styling.CharacterColors{FG: "", BG: ""},
		separatorPrefix: prefix,
		separator:       separator,
		separatorSuffix: suffix,
	}
}

// Segment prints a new Powerline segment with the given background color and text.
func (pl *Powerline) Segment(color string, text string) string {
	// If the segment is empty, skip it.
	if text == "" {
		return ""
	}

	result := ""

	style, err := pl.styles.Get(styling.ToBgColor(color))
	if err != nil {
		return err.Error() + text
	}
	coloredText, firstColor, lastColor := style.ApplyGetColors(text)

	// Print the separator
	if pl.lastColor.BG != "" {
		prefixStyle, err := pl.styles.Get(firstColor.BG + " " + styling.ToBgColor(pl.lastColor.BG))
		if err == nil {
			result += prefixStyle.Apply(pl.separatorPrefix)
		}

		suffixStyle, err := pl.styles.Get(styling.ToBgColor(firstColor.BG) + " " + pl.lastColor.BG)
		if err == nil {
			result += suffixStyle.Apply(pl.separator + pl.separatorSuffix)
		}
	}

	result += coloredText
	pl.lastColor = lastColor
	if pl.lastColor.BG == "" {
		pl.lastColor.BG = "black"
	}

	return result
}

func (pl *Powerline) Finish() string {
	result := ""

	if pl.lastColor.BG != "" {
		prefixStyle, err := pl.styles.Get("black " + styling.ToBgColor(pl.lastColor.BG))
		if err == nil {
			result += prefixStyle.Apply(pl.separatorPrefix)
		}

		suffixStyle, err := pl.styles.Get("bg:black " + pl.lastColor.BG)
		if err == nil {
			result += suffixStyle.Apply(pl.separator)
		}
	}

	return result
}

// TxtFuncMap returns template functions for styling text.
func TxtFuncMap(styles *styling.Registry) template.FuncMap {
	return template.FuncMap{
		"makePowerline": func(prefix string, separator string, suffix string) *Powerline {
			return New(styles, prefix, separator, suffix)
		},
	}
}
