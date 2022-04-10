package ansigradient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	tokens, printWidth := tokenize("abc")
	assert.Equal(t, []gradientToken{
		{t: tokenString, content: "abc", printWidth: 3},
	}, tokens)
	assert.Equal(t, 3, printWidth)
}

func TestTokenizeZwjEmoji(t *testing.T) {
	tokens, printWidth := tokenize("👩🏻‍🚀")
	assert.Equal(t, []gradientToken{
		{t: tokenComplexChar, content: "👩🏻‍🚀", printWidth: 2},
	}, tokens)
	// Astronaut characters is 2 columns wide.
	assert.Equal(t, 2, printWidth)

	tokens, printWidth = tokenize("A👩🏻‍🚀D")
	assert.Equal(t, []gradientToken{
		{t: tokenString, content: "A", printWidth: 1},
		{t: tokenComplexChar, content: "👩🏻‍🚀", printWidth: 2},
		{t: tokenString, content: "D", printWidth: 1},
	}, tokens)
	assert.Equal(t, 4, printWidth)

	tokens, printWidth = tokenize("👩🏻‍🚀D")
	assert.Equal(t, []gradientToken{
		{t: tokenComplexChar, content: "👩🏻‍🚀", printWidth: 2},
		{t: tokenString, content: "D", printWidth: 1},
	}, tokens)
	assert.Equal(t, 3, printWidth)

	tokens, printWidth = tokenize("A👩🏻‍🚀")
	assert.Equal(t, []gradientToken{
		{t: tokenString, content: "A", printWidth: 1},
		{t: tokenComplexChar, content: "👩🏻‍🚀", printWidth: 2},
	}, tokens)
	assert.Equal(t, 3, printWidth)
}
