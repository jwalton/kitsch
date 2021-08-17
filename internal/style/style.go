// Package style is for parsing style strings from the kitsch configuration,
// and for applying styles to strings.
package style

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jwalton/gchalk/pkg/ansistyles"
)

var hexColorRegex = regexp.MustCompile(`([a-fA-F\d]{6}|[a-fA-F\d]{3})`)

const bgPrefix = "bg:"

// Style represents ANSI styling to apply to a string of text in terminal output.
// Style keeps track of a foreground color, a background color, and any addiitional
// modifiers which should be applied to the string.
type Style struct {
	// FG is the foreground color of this style.  This can be any string that
	// `gchalk.Style()` accepts (e.g. "red", "brightBlack") or a hex string (e.g. "#2080ff").
	FG string
	// BG is the background color of this style.
	BG string
	// Modifiers is an array of modifiers (e.g. "bold").  These can be any
	// modifier accepted by `gchalk.Style()`.
	Modifiers []string
}

// IsEmpty returns true if this instance of Style is the default empty instance,
// which applies no styling.
func (style *Style) IsEmpty() bool {
	return style.FG == "" && style.BG == "" && len(style.Modifiers) == 0
}

func (style *Style) reset() {
	style.FG = ""
	style.BG = ""
	style.Modifiers = nil
}

// Mix will mix another style onto this style.  Styles from "nested" will override
// styles from the current style.  This is used in the case where we nest styles -
// this style is the parent style, and "nested" is the nested style.
//
// For example, if this style is `{BG: "red", FG: "white"}` and `nested` is
// `{BG: "blue"}`, this would return `{BG: "blue", FG: "white"}`, which is
// what style you get if you tried to nest content with the "nested" style inside
// this style.
//
func (style *Style) Mix(nested Style) Style {
	return Style{
		FG:        defaultString(nested.FG, style.FG),
		BG:        defaultString(nested.BG, style.BG),
		Modifiers: mixModifiers(style.Modifiers, nested.Modifiers),
	}
}

// Default will return the given style if and only if the receiver is empty.
func (style *Style) Default(defaultStyle Style) Style {
	if style.IsEmpty() {
		return defaultStyle
	}
	return *style
}

// defaultString returns value if it is non-empty, or def otherwise.
func defaultString(value string, def string) string {
	if value != "" {
		return value
	}
	return def
}

func mixModifiers(base []string, nested []string) []string {
	if len(base) == 0 {
		return nested
	}
	if len(nested) == 0 {
		return base
	}

	styles := make(map[string]struct{}, len(base)+len(nested))
	for _, style := range nested {
		styles[style] = struct{}{}
	}
	_, nestedReset := styles["reset"]
	if nestedReset {
		// If there's a reset in the nested styles, just return the nested styles.
		return nested
	}

	for _, style := range base {
		styles[style] = struct{}{}
	}
	_, baseReset := styles["reset"]
	if baseReset {
		// If there's a reset in the base styles, return
		return append(base, nested...)
	}

	result := make([]string, 0, len(styles))
	for style := range styles {
		result = append(result, style)
	}

	return result
}

// Parse parses a style string from the kitsch configuration, and returns
// a `Style` object, which can be passed to `applyStyles()`.
//
// Valid style strings include:
//
// • Any color name accepted by `gchalk.Style()` (e.g. "red", "blue", "brightBlue").
//
// • A hex color code (e.g. "#FFF" or "#320fc9").
//
// • Any of the above, but starting with "bg:" to style the background.
//
// • Any modifier accepted by `gchalk.Style()` (e.g. "bold", "dim", "inverse").
//
// Multiple style strings can be separated by spaces.
//
func Parse(styleString string) (Style, error) {
	result := Style{}
	return result, result.parse(styleString)
}

func (style *Style) parse(styleString string) error {
	var err error

	style.reset()

	substrings := strings.Split(styleString, " ")
	for _, str := range substrings {
		if len(str) > 0 {
			styleErr := parseStyleSubstring(str, false, style)
			if styleErr != nil && err == nil {
				err = styleErr
			}
		}
	}

	return err
}

// ParseMust is like Parse() but panics if the styleString is invalid.
func ParseMust(styleString string) Style {
	result, err := Parse(styleString)

	if err != nil {
		panic(err)
	}

	return result
}

// isBgColor returns true if the passed in style is a background color,
// and returns the color with the bg prefix stripped.
func isBgColor(styleStr string) (string, bool) {
	if strings.HasPrefix(styleStr, bgPrefix) {
		// Handle case where `styleStr` starts with "bg:".
		return styleStr[3:], true
	} else if strings.HasPrefix(styleStr, "bg") {
		return strings.ToLower(styleStr[2:3]) + styleStr[3:], true
	}
	return "", false
}

// parseStyleSubstring is a helper function for `ParseStyle` which parses an individual
// style string.
func parseStyleSubstring(styleStr string, isBackground bool, style *Style) error {
	// TODO: Add `bg#abcdef` support to gchalk.  Maybe add `bg:red` support, too.
	// It would be nice if `gchalk.style()` supported these and FG hex colors.
	if color, isBg := isBgColor(styleStr); isBg {
		// Handle case where `styleStr` starts with "bg:" or "bg".
		return parseStyleSubstring(color, true, style)
	} else if validateColor(styleStr) {
		if isBackground {
			style.BG = styleStr
		} else {
			style.FG = styleStr
		}
	} else if _, ok := ansistyles.Modifier[styleStr]; ok {
		style.Modifiers = append(style.Modifiers, styleStr)
	} else {
		return fmt.Errorf("cannot parse style \"%s\"", styleStr)
	}
	return nil
}

// validateColor retunrs true if the given string if a valid style.FG or style.BG.
func validateColor(color string) bool {
	_, validAnsiStyle := ansistyles.Color[color]
	return validAnsiStyle || hexColorRegex.MatchString(color)
}

// ToFgColor converts a color to a foreground color.
func ToFgColor(color string) string {
	if color == "" {
		color = "black"
	} else if fgColor, isBg := isBgColor(color); isBg {
		color = fgColor
	}
	return color
}

// ToBgColor converts a color to a background color.
func ToBgColor(color string) string {
	if color == "" {
		color = "black"
	} else if !strings.HasPrefix(color, "bg") {
		color = "bg:" + color
	}

	return color
}
