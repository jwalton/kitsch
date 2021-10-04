package powerline

import (
	"text/template"

	"github.com/jwalton/kitsch-prompt/internal/style"
)

// Powerline is a helper object for constructing powerline prompts.
type Powerline struct {
	// color is the current background color.
	color           string
	separatorPrefix string
	separatorSuffix string
}

// New creates a noew Powerline helper object for use in a template.
func New(prefix string, suffix string) *Powerline {
	return &Powerline{
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

	prefixStyle, err := style.Parse(style.ToBgColor(lastColor))
	if err != nil {
		return ""
	}
	suffixStyle, err := style.Parse(lastColor + " " + style.ToBgColor(color))
	if err != nil {
		return ""
	}

	prefix, _, _, _ := prefixStyle.Apply(pl.separatorPrefix)
	suffix, _, _, _ := suffixStyle.Apply(pl.separatorSuffix)

	return prefix + suffix
}

// print prints some text into the current powerline segment.
func (pl *Powerline) print(text string) string {
	if pl.color == "" {
		return text
	}

	style, err := style.Parse(style.ToBgColor(pl.color))
	if err != nil {
		return text
	}

	styledSegment, _, _, err := style.Apply(text)
	if err != nil {
		return text
	}

	return styledSegment
}

// Segment prints a new Powerline segment with the given background color and text.
func (pl *Powerline) Segment(color string, text string) string {
	return pl.updateColor(color) + pl.print(text)
}

// TxtFuncMap returns template functions for styling text.
func TxtFuncMap() template.FuncMap {
	return template.FuncMap{
		"makePowerline": New,
	}
}
