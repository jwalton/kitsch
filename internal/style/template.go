package style

import (
	"fmt"
	"text/template"
)

func toText(text interface{}) string {
	textString, ok := text.(string)
	if !ok {
		textString = fmt.Sprintf("%v", text)
	}
	return textString
}

// Apply a foreground color to the given text.
func fgColor(color string, text interface{}) string {
	textString := toText(text)
	color = ToFgColor(color)

	style := Style{}
	err := style.parse(color)
	if err == nil {
		textString, _, _, _ = style.Apply(textString)
	}
	return textString
}

// BgColor applies a color as a background color to the given text.
func bgColor(color string, text interface{}) string {
	textString := toText(text)
	color = ToBgColor(color)

	style := Style{}
	err := style.parse(color)
	if err == nil {
		textString, _, _, _ = style.Apply(textString)
	}

	return textString
}

// style a string in the given styles.
func style(args ...interface{}) string {
	argc := len(args)
	textString := toText(args[argc-1])
	styles := args[:argc-1]

	styled := textString
	for _, styleDef := range styles {
		style := Style{}
		err := style.UnmarshalInterface(styleDef)
		if err != nil {
			style = Style{}
		}

		styled, _, _, _ = style.Apply(styled)
	}

	return styled
}

// TxtFuncMap returns template functions for styling text.
func TxtFuncMap() template.FuncMap {
	return template.FuncMap{
		"style":   style,
		"fgColor": fgColor,
		"bgColor": bgColor,
	}
}
