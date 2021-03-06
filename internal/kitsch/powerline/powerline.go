package powerline

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jwalton/kitsch/internal/kitsch/styling"
)

// Powerline is a helper object for constructing powerline prompts.
type Powerline struct {
	styles          *styling.Registry
	lastColor       styling.CharacterColors
	separatorPrefix string
	separator       string
	separatorSuffix string
	reverse         bool
}

// New creates a new Powerline helper object for use in a template.
func New(styles *styling.Registry, prefix string, separator string, suffix string, reverse bool) *Powerline {
	return &Powerline{
		styles:          styles,
		lastColor:       styling.CharacterColors{FG: "", BG: ""},
		separatorPrefix: prefix,
		separator:       separator,
		separatorSuffix: suffix,
		reverse:         reverse,
	}
}

// Segment prints a new Powerline segment with the given background color and text.
// If previous segments have been written by this Powerline instance, then
// a "prefix+separator+suffix" will be written between the previous segment and
// this one.
//
// The prefix will have the background from the previous segment and the foreground
// from the next segment.  The suffix will be the reverse.  The separator will be
// colored the same as the suffix, unless this is a "reverse" powerline instance
// in which case it will be the same as the prefix.
func (pl *Powerline) Segment(color string, text interface{}) string {
	// If the segment is empty, skip it.
	if text == nil {
		return ""
	}

	str := toText(text)
	if str == "" {
		return ""
	}

	out := strings.Builder{}

	style, err := pl.styles.Get(styling.ToBgColor(color))
	if err != nil {
		return err.Error() + str
	}
	coloredText, firstColor, lastColor := style.ApplyGetColors(str)

	// Print the separator
	if pl.lastColor.BG != "" || pl.reverse {
		prefix := pl.separatorPrefix
		suffix := pl.separatorSuffix

		if pl.reverse {
			if pl.lastColor.BG == "" {
				prefix = pl.separator
			} else {
				prefix = pl.separatorPrefix + pl.separator
			}
		} else {
			suffix = pl.separator + pl.separatorSuffix
		}

		prefixStyle, err := pl.styles.Get(styling.ToFgColor(firstColor.BG) + " " + styling.ToBgColor(pl.lastColor.BG))
		if err == nil {
			out.WriteString(prefixStyle.Apply(prefix))
		}

		suffixStyle, err := pl.styles.Get(styling.ToFgColor(pl.lastColor.BG) + " " + styling.ToBgColor(firstColor.BG))
		if err == nil {
			out.WriteString(suffixStyle.Apply(suffix))
		}
	}

	out.WriteString(coloredText)
	pl.lastColor = lastColor
	if pl.lastColor.BG == "" {
		pl.lastColor.BG = "bg:black"
	}

	return out.String()
}

// Finish will print an "end" to this powerline string.
func (pl *Powerline) Finish() string {
	out := strings.Builder{}

	if pl.lastColor.BG != "" && !pl.reverse {
		prefixStyle, err := pl.styles.Get("black " + styling.ToBgColor(pl.lastColor.BG))
		if err == nil {
			out.WriteString(prefixStyle.Apply(pl.separatorPrefix))
		}

		suffixStyle, err := pl.styles.Get("bg:black " + styling.ToFgColor(pl.lastColor.BG))
		if err == nil {
			out.WriteString(suffixStyle.Apply(pl.separator))
		}
	}

	return out.String()
}

// TxtFuncMap returns template functions for styling text.
func TxtFuncMap(styles *styling.Registry) template.FuncMap {
	return template.FuncMap{
		"newPowerline": func(prefix string, separator string, suffix string) *Powerline {
			return New(styles, prefix, separator, suffix, false)
		},
		"newReversePowerline": func(prefix string, separator string, suffix string) *Powerline {
			return New(styles, prefix, separator, suffix, true)
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
