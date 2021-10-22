package ansigradient

import (
	"github.com/jwalton/go-ansiparser"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// tokenType represents the type of a parsed token.
type tokenType int

const (
	tokenString      tokenType = 0
	tokenEscapeCode  tokenType = 1
	tokenComplexChar tokenType = 2
)

type gradientToken struct {
	t       tokenType
	content string
	fg      string
	bg      string
	// printWidth is how many columns would the content of this token occupy when printed to the screen.
	printWidth int
}

// tokenize returns the set of tokens, and the sum of their print widths.
func tokenize(s string) ([]gradientToken, int) {
	tokens := []gradientToken{}

	ansiTokenizer := ansiparser.NewStringTokenizer(s)

	for ansiTokenizer.Next() {
		t := ansiTokenizer.Token()
		switch t.Type {
		case ansiparser.String:
			if isASCIIString(t.Content) {
				// If the string is all ASCII, just copy it to a string token.
				tokens = append(tokens, gradientToken{
					t:          tokenString,
					content:    t.Content,
					fg:         t.FG,
					bg:         t.BG,
					printWidth: len(t.Content),
				})
			} else {
				// If there are unicode characters, find them and work out
				// their print widths.
				tokens = append(tokens, tokenizeStringTokenWithUnicodeCharacters(t)...)
			}

		case ansiparser.EscapeCode:
			tokens = append(tokens, gradientToken{
				t:          tokenEscapeCode,
				content:    t.Content,
				fg:         t.FG,
				bg:         t.BG,
				printWidth: 0,
			})
		}
	}

	printWidth := 0
	for _, token := range tokens {
		printWidth += token.printWidth
	}

	return tokens, printWidth
}

// isASCIIString returns true if s contains only ASCII characters.
func isASCIIString(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}

func tokenizeStringTokenWithUnicodeCharacters(ansiToken ansiparser.AnsiToken) []gradientToken {
	tokens := []gradientToken{}
	str := ansiToken.Content
	position := 0

	makeStringToken := func(str string) {
		tokens = append(tokens, gradientToken{
			t:          tokenString,
			content:    str,
			fg:         ansiToken.FG,
			bg:         ansiToken.BG,
			printWidth: len(str),
		})
	}

	// Grab any non-unicode characters at the start of the string.
	for position < len(str) {
		if str[position] > 127 {
			break
		}
		position++
	}
	if position > 0 {
		makeStringToken(str[:position])
	}

	// Handle the rest of the string as unicode.
	str = str[position:]

	position = 0
	start := position
	graphemes := uniseg.NewGraphemes(str)
	for graphemes.Next() {
		grapheme := graphemes.Str()
		if len(grapheme) == 1 {
			position++
		} else {
			// Multi-byte grapheme.

			if position != start {
				// Finish off the string we were working on.
				makeStringToken(str[start:position])
			}

			// Add the grapheme token.
			tokens = append(tokens, gradientToken{
				t:          tokenComplexChar,
				content:    grapheme,
				fg:         ansiToken.FG,
				bg:         ansiToken.BG,
				printWidth: runewidth.StringWidth(grapheme),
			})

			position += len(grapheme)
			start = position
		}
	}

	if start != position {
		makeStringToken(str[start:position])
	}

	return tokens
}
