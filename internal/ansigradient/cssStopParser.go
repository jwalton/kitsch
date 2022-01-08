package ansigradient

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/jwalton/kitsch/internal/colortools"
)

type cssStopParser struct {
	colorMap map[string]string
	str      string
	index    int
}

func parseCSSStops(colorMap map[string]string, cssStops string) ([]gradientStop, error) {
	result := make([]gradientStop, 0, 2)

	parser := cssStopParser{colorMap: colorMap, str: cssStops, index: 0}

	for parser.index < len(parser.str) {
		parser.skipWhitespace()
		segmentStartIndex := parser.index
		currentStopIndex := len(result)

		if parser.index < len(parser.str) && parser.isDigit(parser.str[parser.index]) {
			// Parse the first `length-percentage`, if there is one.
			offset, offsetType, err := parser.parseLengthPercentage()
			if err != nil {
				return result, err
			}

			if currentStopIndex == 0 {
				return result, fmt.Errorf("Expected linear-color-stop at position %d", segmentStartIndex+1)
			} else if result[currentStopIndex-1].ColorUnset {
				return result, fmt.Errorf("Cannot have two linear-color-hint in a row at position %d", segmentStartIndex+1)
			}

			stop := gradientStop{ColorUnset: true, Offset: offset, OffsetType: offsetType}
			result = append(result, stop)
		} else {
			// Parse color.
			color, err := parser.parseColor()
			if err != nil {
				return result, err
			}

			// Parse the first `length-percentage`, if there is one.
			offset, offsetType, err := parser.parseLengthPercentage()
			if err != nil {
				return result, err
			}

			// Parse the second `length-percentage`, if there is one.
			offset2, offsetType2, err := parser.parseLengthPercentage()
			if err != nil {
				return result, err
			}

			stop := gradientStop{Color: color, ColorUnset: false, Offset: offset, OffsetType: offsetType}
			result = append(result, stop)

			// If there are two offsets in the linear-color-stop, then push two stops.
			if offsetType2 != gradientStopUnspecified {
				stop := gradientStop{Color: color, ColorUnset: false, Offset: offset2, OffsetType: offsetType2}
				result = append(result, stop)
			}
		}

		// Find the comma, if there is one.
		parser.skipWhitespace()
		if parser.index < len(parser.str) {
			if parser.str[parser.index] == ',' {
				parser.index++
			} else {
				return result, fmt.Errorf(
					"expected ',' at position %d, found '%v'",
					parser.index+1,
					string(parser.str[parser.index]),
				)
			}
		}
	}

	return result, nil
}

func (parser *cssStopParser) skipWhitespace() {
	// Skip whitepace.
	for parser.index < len(parser.str) && parser.isWhitespace(parser.str[parser.index]) {
		parser.index++
	}
}

func (parser *cssStopParser) isWhitespace(c byte) bool {
	return c == ' ' || c == '\t'
}

func (parser *cssStopParser) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (parser *cssStopParser) getNextToken() string {
	str := parser.str
	tokenStart := parser.index

	for parser.index < len(str) && !parser.isWhitespace(str[parser.index]) && str[parser.index] != ',' {
		parser.index++
	}

	return str[tokenStart:parser.index]
}

func (parser *cssStopParser) parseColor() (c color.RGBA, err error) {
	// TODO: Add support for `rgb()`, for named colors, and for a color dictionary of custom colors.
	colorStartIndex := parser.index

	token := parser.getNextToken()

	if token[0] == '$' {
		// Try for a custom color.
		if colorStr, ok := parser.colorMap[token]; ok {
			c, err = colortools.ParseColor(colorStr)
			if err != nil {
				return c, fmt.Errorf("color %s=\"%s\" cannot be used in linear-gradient", token, colorStr)
			}
			return c, nil
		}
	}

	c, err = colortools.ParseColor(token)
	if err != nil {
		return c, fmt.Errorf("invalid color \"%s\" at %d", token, colorStartIndex)
	}

	return c, nil
}

// parseLengthPercentage parses a `length-percentage` as per CSS values level 4:
// https://www.w3.org/TR/css-values-4/#typedef-length-percentage
func (parser *cssStopParser) parseLengthPercentage() (offset float64, offsetType gradientOffsetType, err error) {
	str := parser.str

	parser.skipWhitespace()

	offsetStartIndex := parser.index

	if parser.index >= len(str) || !parser.isDigit(str[parser.index]) {
		// No offset
		return 0, gradientStopUnspecified, nil
	}

	start := parser.index
	for parser.index < len(str) && parser.isDigit(str[parser.index]) {
		parser.index++
	}

	if parser.index <= len(str) {
		offsetStr := str[start:parser.index]
		offset, err = strconv.ParseFloat(offsetStr, 64)
		if err != nil {
			return 0, gradientStopUnspecified, err
		}
	}

	if parser.index < len(str) && str[parser.index] == '%' {
		// Percent offset
		parser.index++
		return offset / 100, gradientStopRelative, nil
	} else if parser.index < (len(str)-1) && str[parser.index] == 'p' && str[parser.index+1] == 'x' {
		// Pixel offset
		parser.index += 2
		return offset, gradientStopAbsolute, nil
	} else if offset == 0 {
		// Treat bare "0" as a pixel offset.
		return 0, gradientStopAbsolute, nil
	} else {
		return 0, gradientStopUnspecified, fmt.Errorf("Expected offset at position %d, got %v", offsetStartIndex+1, str[parser.index:])
	}
}
