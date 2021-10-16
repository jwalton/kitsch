package powerline

import (
	"text/template"

	"github.com/jwalton/kitsch-prompt/internal/style"
)

// Powerline is a helper object for constructing powerline prompts.
type Powerline struct {
	styles *style.Registry
	// color is the current background color.
	color           string
	separatorPrefix string
	separatorSuffix string
}

// New creates a noew Powerline helper object for use in a template.
func New(styles *style.Registry, prefix string, suffix string) *Powerline {
	return &Powerline{
		styles:          styles,
		color:           "",
		separatorPrefix: prefix,
		separatorSuffix: suffix,
	}
}

// updateColor updates the current color of this powerline object, and
// returns a powerline separator.
func (pl *Powerline) updateColor(color string) string {
	if color == "" {
		color = "black"
	}

	lastColor := pl.color
	pl.color = style.ToFgColor(color)

	// First segment and segments where color doesn't change get no separator.
	if lastColor == "" || lastColor == color {
		return ""
	}

	prefixStyle, err := pl.styles.Get(style.ToBgColor(lastColor))
	if err != nil {
		return ""
	}
	suffixStyle, err := pl.styles.Get(lastColor + " " + style.ToBgColor(color))
	if err != nil {
		return ""
	}

	prefix := prefixStyle.Apply(pl.separatorPrefix)
	suffix := suffixStyle.Apply(pl.separatorSuffix)

	return prefix + suffix
}

// print prints some text into the current powerline segment.
func (pl *Powerline) print(text string) string {
	if pl.color == "" {
		return text
	}

	style, err := pl.styles.Get(style.ToBgColor(pl.color))
	if err != nil {
		return text
	}

	styledSegment := style.Apply(text)

	return styledSegment
}

// Segment prints a new Powerline segment with the given background color and text.
func (pl *Powerline) Segment(color string, text string) string {
	return pl.updateColor(color) + pl.print(text)
}

// TxtFuncMap returns template functions for styling text.
func TxtFuncMap(styles *style.Registry) template.FuncMap {
	return template.FuncMap{
		"makePowerline": func(prefix string, suffix string) func(string, string) string {
			return func(color string, text string) string {
				return New(styles, prefix, suffix).Segment(color, text)
			}
		},
	}
}
