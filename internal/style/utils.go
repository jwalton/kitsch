package style

import "strings"

const bgPrefix = "bg:"

// isBgColor returns true if the passed in style is a background color,
// and returns the color with the bg prefix stripped.
func isBgColor(styleStr string) (string, bool) {
	if strings.HasPrefix(styleStr, bgPrefix) {
		// Handle case where `styleStr` starts with "bg:".
		return styleStr[3:], true
	} else if strings.HasPrefix(styleStr, "bg") {
		// Handle bgRed or bg#ff0000 style colors.
		return strings.ToLower(styleStr[2:3]) + styleStr[3:], true
	}
	return "", false
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
