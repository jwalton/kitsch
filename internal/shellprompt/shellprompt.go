// Package shellprompt provides functions for outputting the prompt for different kinds
// of shells.
package shellprompt

import (
	"github.com/jwalton/go-ansiparser"
)

// AddZeroWidthCharacterEscapes adds special characters around zero length strings
// to the provided prompt.
func AddZeroWidthCharacterEscapes(shell string, prompt string) string {
	switch shell {
	case "zsh":
		// https: //zsh.sourceforge.io/Doc/Release/Prompt-Expansion.html#Visual-effects
		return addZeroWidthCharacterEscapes(prompt, "%{", "%}")
	case "bash":
		// https://www.gnu.org/software/bash/manual/html_node/Controlling-the-Prompt.html#Controlling-the-Prompt
		return addZeroWidthCharacterEscapes(prompt, "\\[", "\\]")
	}

	return prompt
}

func addZeroWidthCharacterEscapes(prompt string, start string, end string) string {
	parsed := ansiparser.Parse(prompt)
	result := ""

	for _, part := range parsed {
		if part.Type == ansiparser.EscapeCode {
			result += start + part.Content + end
		} else {
			result += part.Content
		}
	}

	return result
}
