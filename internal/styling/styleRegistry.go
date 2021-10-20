package styling

import (
	"fmt"

	"github.com/jwalton/gchalk"
)

// Registry is used to store and retrieve styles.
type Registry struct {
	// CustomColors is a map of color names and their replacements.  For example,
	// if CustomColors["$foregroud"] = "red", then "$foreground" could be used in
	// a style string to refer to the color red.  Custom colors must start with
	// a "$".
	CustomColors   map[string]string
	styles         map[string]*Style
	gchalkInstance *gchalk.Builder
}

// AddCustomColor registers a custom color with the registry.
// The color name must start with a "$".
func (registry *Registry) AddCustomColor(name string, color string) {
	if registry.CustomColors == nil {
		registry.CustomColors = map[string]string{}
	}
	registry.CustomColors[name] = color
}

// Get compiles a style string into a style, and returns the style.  Styles
// are cached in the registry, so getting the same styleString twice will return
// the same Style object.
//
// Valid style strings include any of the following, separated by spaces::
//
// • Any color name accepted by `gchalk.Style()` (e.g. "red", "blue", "brightBlue").
//
// • A hex color code (e.g. "#FFF" or "#320fc9").
//
// • A CSS style linear-gradient (e.g. "linear-gradient(#f00, #00f)".
//
// • A custom color (e.g. "$foreground").
//
// • Any of the above, but starting with "bg:" to style the background.
//
// • Any modifier accepted by `gchalk.Style()` (e.g. "bold", "dim", "inverse").
//
func (registry *Registry) Get(styleString string) (*Style, error) {
	if style := registry.styles[styleString]; style != nil {
		return style, nil
	}

	// Lazy initialization of the registry.
	if registry.gchalkInstance == nil {
		builder, err := gchalk.WithStyle()
		if err != nil {
			return nil, err
		}
		registry.gchalkInstance = builder
	}

	if registry.styles == nil {
		registry.styles = map[string]*Style{}
	}

	style, err := compileStyle(registry.gchalkInstance, registry.CustomColors, styleString)
	if err != nil {
		return nil, fmt.Errorf("error compiling style \"%s\": %w", styleString, err)
	}

	registry.styles[styleString] = &style

	return &style, nil
}
