// Package ansigradient is used to apply gradients and colors to strings using
// ANSI escape codes.
//
// This library can be used to create gradients, and then to apply those
// gradients (or any arbitrary collection of colors) to a string.
//
// To use this library, first you create a `Gradient` using the same syntax as
// used in CSS `linear-gradient()` (except that you can't specify a direction -
// they are always left-to-right).  The simplest gradient would have just two stops:
//
//     // Create a gradient.
//     linearGradient := ansigradient.CSSLinearGradientMust("0xff0000, 0x0000ff")
//
//     // Apply a gradient to the foreground of a string of text.
//     text := ansigradient.ApplyGradients("Red-to-blue text", linearGradient, nil);
//
//     // Apply a gradient to the background of a string of text.
//     text := ansigradient.ApplyGradients("Red-to-blue text", nil, linearGradient);
//
// Stops can optionally have an "offset" which will say where along the text that
// color will begin - an offset of `0px` or `0%` is right before the first color,
// and `100%` is right at the end of the last character.  The actual color of
// any character is computed as if the character were one pixel wide and is
// based on the "center" of that character (i.e. the first character gets it's
// color from linearly-interpolating the color at 0.5px).  As with CSS, we can
// also provide a stop with no color - this is used as a hint for where the
// midpoint between two colors should be:
//
//     gradient := ansigradient.CSSLinearGradientMust("#FF0000 20%, 30%, #0000FF 80%")
//
// The "ApplyGradients"  functions will attempt to auto-detect terminal color
// support based on what stdout supports.  For details about how
// this works, see https://github.com/jwalton/go-supportscolor.  You can
// override the level with `SetLevel()`, or if you want to apply colors to a
// string ignoring the current color level, you can do so with
// "ApplyGradientsRaw".
//
package ansigradient

import (
	"image/color"
)

type gradientOffsetType int

const (
	gradientStopUnspecified gradientOffsetType = 0
	gradientStopAbsolute    gradientOffsetType = 1
	gradientStopRelative    gradientOffsetType = 2
)

type gradientStop struct {
	// Color of this stop.
	Color color.RGBA
	// ColorUnset is true if the color for this stop has not been set.
	ColorUnset bool
	// Offset is the offset of this stop.  This is either a value in pixels if
	// Absolute is true, or a value between 0 and 1 if relative.  If the value
	// is less than 0, it is treated as unspecified.
	Offset float64
	// OffsetType is one of gradientStopUnspecified, gradientStopAbsolute, or gradientStopRelative.
	OffsetType gradientOffsetType
}

// HasUndefinedOffset returns true if the offset is undefined.
func (stop *gradientStop) HasUndefinedOffset() bool {
	return stop.OffsetType == gradientStopUnspecified
}

// GetOffset returns the absolute offset of a stop, given the length of the
// gradient being generated.
func (stop *gradientStop) GetOffset(length int) float64 {
	switch stop.OffsetType {
	case gradientStopUnspecified:
		panic("cannot get offset from stop with undefined offset")
	case gradientStopAbsolute:
		return stop.Offset
	case gradientStopRelative:
		return float64(length) * stop.Offset
	default:
		panic("unknown offset type")
	}
}

// Gradient represents a color gradient, or any object which can generate a series of colors.
type Gradient interface {
	// Colors returns an array of colors for this gradient, of the specified length.
	Colors(length int) []color.RGBA
	// ColorAt returns the color of a specific location along the gradient.
	ColorAt(length int, index int) color.RGBA
	// Generator returns a Colorizer for the specified length string.
	Generator(length int) ColorGenerator
}

// ColorGenerator generates colors along a gradient or spectrum.
type ColorGenerator interface {
	// ColorAt returns the color of a specific location along a color spectrum,
	// such as a gradient.  `position` is the position along the spectrum.
	// The range of `position` depends on the length of the spectrum.
	ColorAt(position float64) color.RGBA
}
