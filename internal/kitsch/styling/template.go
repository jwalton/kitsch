package styling

import (
	"fmt"
	"text/template"
)

// toText converts a value to a string.
func toText(text interface{}) string {
	textString, ok := text.(string)
	if !ok {
		textString = fmt.Sprintf("%v", text)
	}
	return textString
}

// TxtFuncMap returns template functions for styling text.
func TxtFuncMap(styles *Registry) template.FuncMap {
	// fgColor applies a foreground color to the given text.
	fgColor := func(color string, text interface{}) string {
		textString := toText(text)
		color = ToFgColor(color)

		style, err := styles.Get(color)
		if err == nil {
			textString = style.Apply(textString)
		}
		return textString
	}

	// bgColor applies a color as a background color to the given text.
	bgColor := func(color string, text interface{}) string {
		textString := toText(text)
		color = ToBgColor(color)

		style, err := styles.Get(color)
		if err == nil {
			textString = style.Apply(textString)
		}
		return textString
	}

	// style a string in the given style.
	style := func(styleStr string, text interface{}) string {
		styled := toText(text)

		style, err := styles.Get(styleStr)
		if err == nil {
			styled = style.Apply(styled)
		}

		return styled
	}

	return template.FuncMap{
		"style":   style,
		"fgColor": fgColor,
		"bgColor": bgColor,
	}
}
