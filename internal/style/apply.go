package style

import (
	"strings"

	"github.com/jwalton/gchalk"
)

// Apply is called to apply this style to a string of text.
func (style Style) Apply(text string) (
	result string,
	firstStyle Style,
	lastStyle Style,
	err error,
) {
	if style.IsEmpty() {
		return text, Style{}, Style{}, nil
	}

	text, colorErr := applyColor(text, style.FG, false)
	if colorErr != nil {
		err = colorErr
	}
	text, colorErr = applyColor(text, style.BG, true)
	if colorErr != nil {
		err = colorErr
	}
	mods, colorErr := gchalk.Style(style.Modifiers...)
	if colorErr != nil {
		err = colorErr
	} else {
		text = mods(text)
	}

	// TODO: When we add gradient support, need to return different styles here.
	return text, style, style, err
}

func applyColor(text string, styleStr string, background bool) (string, error) {
	if styleStr == "" {
		return text, nil
	}

	var colorFn func(text ...string) string
	var err error = nil

	if strings.HasPrefix(styleStr, "#") {
		if background {
			colorFn = gchalk.BgHex(styleStr)
		} else {
			colorFn = gchalk.Hex(styleStr)
		}
	} else {
		if background {
			styleStr = "bg" + strings.ToUpper(styleStr[0:1]) + styleStr[1:]
			colorFn, err = gchalk.Style(styleStr)
		} else {
			colorFn, err = gchalk.Style(styleStr)
		}
	}

	return colorFn(text), err

}
