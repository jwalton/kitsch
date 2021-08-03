package style

import (
	"fmt"

	"github.com/jwalton/gchalk/pkg/ansistyles"
	"gopkg.in/yaml.v3"
)

// UnmarshalYAML will convert a YAML string or object into a Style.
func (style *Style) UnmarshalYAML(node *yaml.Node) error {
	if node == nil {
		return fmt.Errorf("No value provided")
	}

	switch node.Kind {
	case yaml.ScalarNode:
		return style.UnmarshalInterface(node.Value)
	case yaml.MappingNode:
		result := map[string]interface{}{}
		err := node.Decode(&result)
		if err != nil {
			return err
		}
		return style.UnmarshalInterface(result)
	default:
		return fmt.Errorf("Cannot convert node of type %v to style (%d:%d)", node.Kind, node.Line, node.Column)
	}
}

// UnmarshalInterface will parse style from a string, or from a interface{} obtained
// from parsing YAML, JSON, etc...  This will accept either a string or a map
// of the shape `{fg: string, bg: string, modifiers: []string}`.
func (style *Style) UnmarshalInterface(styleInterface interface{}) error {
	// TODO: Can maybe simplify this by just unmarshalling from YAML directly?
	switch styleInterface := styleInterface.(type) {
	case Style:
		style.BG = styleInterface.BG
		style.FG = styleInterface.FG
		style.Modifiers = styleInterface.Modifiers
		return nil
	case *Style:
		style.BG = styleInterface.BG
		style.FG = styleInterface.FG
		style.Modifiers = styleInterface.Modifiers
		return nil
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
					imods, ok := value.([]interface{})
					if !ok {
						return fmt.Errorf("Expected array of strings for style.%v, got: %T - %v", key, mods, mods)
					}
					mods = make([]string, len(imods))
					for i, mod := range imods {
						mods[i] = mod.(string)
					}
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
