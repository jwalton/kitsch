package style

import (
	"fmt"
	"strings"

	"github.com/jwalton/gchalk/pkg/ansistyles"
)

const linearGradientPrefix = "linear-gradient("

type styleDescriptor struct {
	// fg is the foreground color of this style.  This can be any string that
	// `gchalk.Style()` accepts (e.g. "red", "brightBlack"), a hex string
	// (e.g. "#2080ff"), or a CSS style linear-gradient.
	fg string
	// bg is the background color of this style.
	bg string
	// modifiers is an array of modifiers (e.g. "bold").  These can be any
	// modifier accepted by `gchalk.Style()`.
	modifiers []string
}

// parseStyle converts a style string into a style descriptor.
func parseStyle(
	customColors map[string]string,
	styleString string,
) (descriptor styleDescriptor, err error) {
	parser := styleParser{styleString: styleString, position: 0}

	done := false
	for !done {
		token, err := parser.nextToken()
		if err != nil {
			return descriptor, err
		}

		if token == "" {
			done = true
		} else {
			err := parseStyleToken(customColors, &descriptor, token)
			if err != nil {
				return descriptor, err
			}
		}
	}

	return descriptor, nil
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

type styleParser struct {
	styleString string
	position    int
}

// nextToken returns the next token from the style string, if there is one.
// Returns ("", nil) if the end of the string has been reached.
func (parser *styleParser) nextToken() (token string, err error) {
	// Skip any whitespace.
	for parser.position < len(parser.styleString) && isWhitespace(parser.styleString[parser.position]) {
		parser.position++
	}

	// If we're at the end, return io.EOF.
	if parser.position >= len(parser.styleString) {
		return "", nil
	}

	// Read the token.
	tokenStart := parser.position
	for parser.position < len(parser.styleString) && !isWhitespace(parser.styleString[parser.position]) {
		parser.position++

		// If there's a "(", read until the matching ")" and include it in the token.
		if parser.position < len(parser.styleString) && parser.styleString[parser.position] == '(' {
			parser.position++
			for parser.position < len(parser.styleString) && parser.styleString[parser.position] != ')' {
				parser.position++
			}

			if parser.position >= len(parser.styleString) {
				return "", fmt.Errorf("Token at position %v is missing ')'", tokenStart+1)
			}
		}
	}

	return parser.styleString[tokenStart:parser.position], nil
}

// isColor returns true if the given string is a color - eiter .
func isColor(color string) bool {
	_, validAnsiStyle := ansistyles.Color[color]
	if validAnsiStyle {
		return true
	}

	hexColor := parseHexColor(color)
	if hexColor != "" {
		return true
	}

	if strings.HasPrefix(color, linearGradientPrefix) {
		return true
	}

	return false
}

func isHexDigit(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

func parseHexColor(str string) string {
	index := 0

	// Find the "#"
	if index < len(str) && str[index] == '#' {
		index++
	}

	hexStart := index
	for index < len(str) && isHexDigit(str[index]) {
		index++
	}

	colorStr := str[hexStart:index]
	if len(colorStr) != 3 && len(colorStr) != 6 {
		return ""
	}

	return colorStr
}

// parseStyleToken adds a token to the style descriptor.  The token can be
// a color, a modifier, or a linear-gradient.
func parseStyleToken(
	customColors map[string]string,
	descriptor *styleDescriptor,
	token string,
) error {
	return parseStyleTokenHelper(customColors, descriptor, token, false)
}

// parseStyleSubstring is a helper function for `ParseStyle` which parses an individual
// style string.
func parseStyleTokenHelper(
	customColors map[string]string,
	descriptor *styleDescriptor,
	token string,
	isBackground bool,
) error {
	if color, isBg := isBgColor(token); isBg {
		// Handle case where `token` starts with "bg:" or "bg".
		err := parseStyleTokenHelper(customColors, descriptor, color, true)
		if err != nil {
			return fmt.Errorf("unknown style \"%s\"", token)
		}
	} else if isColor(token) {
		// Handle case where `token` is a color.
		if isBackground {
			descriptor.bg = token
		} else {
			descriptor.fg = token
		}
	} else if _, ok := ansistyles.Modifier[token]; ok {
		// Handle case where `token` is a modifier.
		descriptor.modifiers = append(descriptor.modifiers, token)
	} else if color, ok := customColors[token]; ok {
		// Handle casae where token is a custom color.
		if isBackground {
			descriptor.bg = color
		} else {
			descriptor.fg = color
		}
	} else {
		return fmt.Errorf("unknown style \"%s\"", token)
	}

	return nil
}
