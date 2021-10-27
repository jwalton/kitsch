package colortools

import (
	"fmt"
	"image/color"

	"github.com/jwalton/gchalk/pkg/ansistyles"
)

var errUnrecognized = fmt.Errorf("unrecognized color")

func isHexDigit(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

// ValidateHexColor returns true if str is a valid hex color, false otherwise.
func ValidateHexColor(str string) bool {
	index := 0

	// Find the "#"
	if index < len(str) && str[index] == '#' {
		index++
	} else {
		return false
	}

	for index < len(str) && isHexDigit(str[index]) {
		index++
	}

	if index != 4 && index != 7 {
		return false
	}

	return true
}

// HexToColor converts a hex value to a color.
func HexToColor(str string) (color.RGBA, error) {
	if !ValidateHexColor(str) {
		return color.RGBA{}, fmt.Errorf("invalid hex color")
	}

	r, g, b := ansistyles.HexToRGB(str)
	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

// ColorToHex converts a color to a hex string.
func ColorToHex(c color.RGBA) string {
	// TODO: Can make this much more performant.
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// ParseColor converts a string to an RGBA color.
func ParseColor(str string) (color.RGBA, error) {
	if str == "" {
		return color.RGBA{}, errUnrecognized
	} else if str[0] == '#' {
		return HexToColor(str)
	} else if c, ok := CSSColors[str]; ok {
		return c, nil
	}

	return color.RGBA{}, errUnrecognized
}

// ValidateColor returns true if str is a valid color, false otherwise.
func ValidateColor(str string) bool {
	if str == "" {
		return false
	} else if str[0] == '#' {
		return ValidateHexColor(str)
	} else if _, ok := CSSColors[str]; ok {
		return true
	}

	return false
}
