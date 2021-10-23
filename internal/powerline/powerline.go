package powerline

import (
	"fmt"
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
// If previous segments have been written by this Powerline instance, then
// a "prefix+separator+suffix" will be written between the previous segment and
// this one.
func (pl *Powerline) Segment(color string, text interface{}) string {
	// If the segment is empty, skip it.
	if text == nil {
		return ""
	}

	str := toText(text)
	if str == "" {
		return ""
	}

	result := ""

	style, err := pl.styles.Get(styling.ToBgColor(color))
	if err != nil {
		return err.Error() + str
	}
	coloredText, firstColor, lastColor := style.ApplyGetColors(str)

	// Print the separator
	if pl.lastColor.BG != "" {
		prefixStyle, err := pl.styles.Get(styling.ToFgColor(firstColor.BG) + " " + styling.ToBgColor(pl.lastColor.BG))
		if err == nil {
			result += prefixStyle.Apply(pl.separatorPrefix)
		}

		suffixStyle, err := pl.styles.Get(styling.ToFgColor(pl.lastColor.BG) + " " + styling.ToBgColor(firstColor.BG))
		if err == nil {
			result += suffixStyle.Apply(pl.separator + pl.separatorSuffix)
		}
	}

	result += coloredText
	pl.lastColor = lastColor
	if pl.lastColor.BG == "" {
		pl.lastColor.BG = "bg:black"
	}

	return result
}

// Finish will print an "end" to this powerline string.
func (pl *Powerline) Finish() string {
	result := ""

	if pl.lastColor.BG != "" {
		prefixStyle, err := pl.styles.Get("black " + styling.ToBgColor(pl.lastColor.BG))
		if err == nil {
			result += prefixStyle.Apply(pl.separatorPrefix)
		}

		suffixStyle, err := pl.styles.Get("bg:black " + styling.ToFgColor(pl.lastColor.BG))
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

// toText converts a value to a string.
func toText(text interface{}) string {
	textString, ok := text.(string)
	if !ok {
		textString = fmt.Sprintf("%v", text)
	}
	return textString
}
