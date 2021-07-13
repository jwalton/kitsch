package style

import (
	"fmt"

	"github.com/jwalton/gchalk/pkg/ansistyles"
)

// UnmarshalInterface will parse style from a string, or from a interface{} obtained
// from parsing YAML, JSON, etc...  This will accept either a string or a map
// of the shape `{fg: string, bg: string, modifiers: []string}`.
func (style *Style) UnmarshalInterface(styleInterface interface{}) error {
	switch styleInterface := styleInterface.(type) {
	case string:
		return style.parse(styleInterface)
	case map[string]interface{}:
		style.reset()

		validateColorInterface := func(key string, color interface{}) (string, error) {
			colorStr, ok := color.(string)
			if !ok {
				return "", fmt.Errorf("Expected string for style.%v, got: %T - %v", key, color, color)
			}
			if !validateColor(colorStr) {
				return "", fmt.Errorf("Invalid style.%v: %v", key, colorStr)
			}
			return colorStr, nil
		}

		for key, value := range styleInterface {
			switch key {
			case "fg", "FG", "Fg":
				fg, err := validateColorInterface(key, value)
				if err != nil {
					return err
				}
				style.FG = fg
			case "bg", "BG", "Bg":
				bg, err := validateColorInterface(key, value)
				if err != nil {
					return err
				}
				style.BG = bg
			case "modifiers", "Modifiers":
				mods, ok := value.([]string)
				if !ok {
					return fmt.Errorf("Expected array of strings for style.%v, got: %T - %v", key, mods, mods)
				}
				for _, mod := range mods {
					if _, ok := ansistyles.Modifier[mod]; !ok {
						return fmt.Errorf("Unknown style modifier: %v", mod)
					}
				}
				style.Modifiers = mods
			}
		}
		return nil
	default:
		return fmt.Errorf("Don't know how to parse style of type: %T - %v", styleInterface, styleInterface)
	}
}
