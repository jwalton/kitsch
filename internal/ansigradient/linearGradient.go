package ansigradient

import (
	"fmt"
	"image/color"
)

// LinearGradient represents a linear gradient.
type LinearGradient struct {
	stops []gradientStop
}

// CSSLinearGradient constructs a linear gradient from CSS stops.
//
// Each stop should be a string specifying either a color followed by one or
// two stop positions (specified either as a percentage or a in pixels), or am
// interpolation hint defining how the gradient progresses between adjacent color
// stops.
//
// Because this creates a 2D gradient, angles and sides or corners are not allowed.
// HTML color names are also not allowed, only hex color codes.
//
// The following example would create a gradient that goes from red to blue:
//
//     CSSLinearGradient("0xf00", "0x00f")
//
// The following would create a gradient that is solid red for the first 10 pixels,
// then transitions to blue with the midpoint being 30% along the distance
// between them:
//
//	   CSSLinearGradient("0xf00 0 10px", "30%", "0x00f")
//
// You can pass an array of stops, or a single string with stops separated by commas:
//
//	   CSSLinearGradient("0xf00 0 10px, 30%, 0x00f")
//
func CSSLinearGradient(stops string) (Gradient, error) {
	gradientStops, err := parseCSSStops(nil, stops)
	if err != nil {
		return nil, err
	}

	result := LinearGradient{
		stops: gradientStops,
	}

	if len(result.stops) == 0 {
		return nil, fmt.Errorf("Can't create CSSLinearGradient with no stops")
	}

	result.cleanStops()

	return result, nil
}

// CSSLinearGradientWithMap creates a gradient, and allows supplying a map of custom colors.
func CSSLinearGradientWithMap(colorMap map[string]string, stops string) (Gradient, error) {
	gradientStops, err := parseCSSStops(colorMap, stops)
	if err != nil {
		return nil, err
	}

	result := LinearGradient{
		stops: gradientStops,
	}

	if len(result.stops) == 0 {
		return nil, fmt.Errorf("Can't create CSSLinearGradient with no stops")
	}

	result.cleanStops()

	return result, nil
}

// CSSLinearGradientMust is a convenince function for creating a gradient from CSS
// stops.  Calling this is equivalent to calling CSSLinearGradient, but will
// panic on a parsing error.
func CSSLinearGradientMust(stops string) Gradient {
	gradient, err := CSSLinearGradient(stops)
	if err != nil {
		panic(err)
	}
	return gradient
}

// cleanStops will sanitize the stops defined in a gradient.  This will ensure
// that the first and last stop have a defined color and offset, and will
// interpolate any missing colors, and - if possible - missing stop offsets.
func (gradient *LinearGradient) cleanStops() {

	// If the first/last stops do not have a specified offset, set them to 0% and 100%.
	if gradient.stops[0].HasUndefinedOffset() {
		gradient.stops[0].Offset = 0
		gradient.stops[0].OffsetType = gradientStopRelative
	}
	lastIndex := len(gradient.stops) - 1
	if gradient.stops[lastIndex].HasUndefinedOffset() {
		gradient.stops[lastIndex].Offset = 1
		gradient.stops[lastIndex].OffsetType = gradientStopRelative
	}

	// First and last stops need to have a color associated with them.  If none
	// are specified, we pick black for the first and white for the last.
	if gradient.stops[0].ColorUnset {
		gradient.stops[0].Color = color.RGBA{0, 0, 0, 255}
	}
	lastIndex = len(gradient.stops) - 1
	if gradient.stops[lastIndex].ColorUnset {
		gradient.stops[lastIndex].Color = color.RGBA{255, 255, 255, 255}
	}

	// Interpolate any stops that don't have an offset.
	gradient.interpolateStopOffsets(gradient.stops)

	// Interploate any missing colors.  CSS's linear-gradient only lets you have
	// one stop with no color, which is used to set the "midpoint" between two
	// colors.  We let you specify multiple stops with no color - if you specify two,
	// you're specifying the "33%" and the "66%" points.
	lastDefinedColorIndex := 0
	for i := 1; i < len(gradient.stops); i++ {
		if !gradient.stops[i].ColorUnset {
			if lastDefinedColorIndex != i-1 {
				missingColorCount := float64(i - lastDefinedColorIndex - 1)
				for missingIndex := lastDefinedColorIndex + 1; missingIndex < i; missingIndex++ {
					gradient.stops[missingIndex].Color = lerpColor(
						gradient.stops[lastDefinedColorIndex].Color,
						gradient.stops[i].Color,
						float64(missingIndex-lastDefinedColorIndex)/(missingColorCount+1),
					)
				}
			}
			lastDefinedColorIndex = i
		}
	}
}

// interpolateStopOffsets will, given an array of stops where some stops have
// unspecified offsets, set the missing offsets.
//
// If all missing offsets are between two "absolute" stops or between two relative
// stops, then we need no extra information for complete this.  If not, then since
// we don't yet know the length of the gradient, we will lave some stops
// undefined.
//
func (gradient *LinearGradient) interpolateStopOffsets(stops []gradientStop) {
	lastDefinedStopIndex := 0
	for i := 1; i < len(stops); i++ {
		if !stops[i].HasUndefinedOffset() {
			if lastDefinedStopIndex != i-1 {
				lastIsAbsolute := stops[lastDefinedStopIndex].OffsetType == gradientStopAbsolute
				thisIsAbsolute := stops[i].OffsetType == gradientStopAbsolute

				// Can treat a "0" offset as absolute or relative.
				if stops[lastDefinedStopIndex].Offset == 0 {
					lastIsAbsolute = thisIsAbsolute
				}

				if lastIsAbsolute == thisIsAbsolute {
					interpolationType := stops[i].OffsetType

					// Interpolate the missing stops.
					interpolatedStops := linearInterpolateFloat64(
						stops[lastDefinedStopIndex].Offset,
						stops[i].Offset, i-lastDefinedStopIndex+1,
					)

					for index, stop := range interpolatedStops {
						stops[lastDefinedStopIndex+index].Offset = stop
						stops[lastDefinedStopIndex+index].OffsetType = interpolationType
					}
				}
			}
			lastDefinedStopIndex = i
		}
	}
}

func getStopOffset(stops []gradientStop, stopIndex int, length int) float64 {
	if stopIndex >= len(stops) {
		return float64(length)
	}

	if stops[stopIndex].HasUndefinedOffset() {
		if stopIndex == 0 {
			return 0
		} else if stopIndex == (len(stops) - 1) {
			return float64(length)
		}

		prevOffset := stops[stopIndex-1].GetOffset(length)
		nextOffset := stops[stopIndex+1].GetOffset(length)

		if nextOffset <= prevOffset {
			return prevOffset
		}

		return (nextOffset-prevOffset)/2 + prevOffset
	}

	return stops[stopIndex].GetOffset(length)
}

// ColorAt will, given the length of a string, return a function that will
// return the color at any point along the length of that string.
func (gradient LinearGradient) ColorAt(length int, index int) color.RGBA {
	if index < 0 {
		return gradient.stops[0].Color
	}

	// Keep track of which stop we're on.  This works so long as we're moving
	// forward along the gradient.  If we move backwards, we have to start over
	// from the beginning.
	currentStop := 0
	currentStopOffset := getStopOffset(gradient.stops, 0, length)
	nextStopOffset := getStopOffset(gradient.stops, 1, length)

	// Add 0.5, because we want the middle of the pixel.
	position := float64(index) + 0.5

	for nextStopOffset <= position && currentStop < len(gradient.stops)-1 {
		// Advance to the next stop
		currentStop++
		currentStopOffset = nextStopOffset
		nextStopOffset = gradient.stops[currentStop].GetOffset(length)
	}

	// Calculate the color at this point.
	if currentStop >= len(gradient.stops)-1 {
		return gradient.stops[currentStop].Color
	}

	return lerpColor(
		gradient.stops[currentStop].Color,
		gradient.stops[currentStop+1].Color,
		(position-currentStopOffset)/float64(nextStopOffset-currentStopOffset),
	)

}

// Colors returns an array of colors of the specified length.
//
// This assumes that each returned value is a single pixel.  The color
// returned is for the center of the pixel.
func (gradient LinearGradient) Colors(length int) []color.RGBA {
	result := make([]color.RGBA, length)

	// Render all the stops.
	currentStop := 0
	currentStopOffset := getStopOffset(gradient.stops, 0, length)
	nextStopOffset := getStopOffset(gradient.stops, 1, length)

	for i := 0; i < length; i++ {
		// Add 0.5, because we want the middle of the pixel.
		position := float64(i) + 0.5

		// Advance the stop until the position is after the start of the current stop.
		for nextStopOffset <= position && currentStop < len(gradient.stops)-1 {
			currentStop++
			currentStopOffset = nextStopOffset
			nextStopOffset = getStopOffset(gradient.stops, currentStop+1, length)
		}

		if currentStop == 0 && position < currentStopOffset {
			// If we're before the first stop, use the first stop as the color.
			result[i] = gradient.stops[0].Color
		} else if currentStop >= len(gradient.stops)-1 {
			// If we're at or after the last stop, use the last stop as the color.
			result[i] = gradient.stops[currentStop].Color
		} else {
			result[i] = lerpColor(
				gradient.stops[currentStop].Color,
				gradient.stops[currentStop+1].Color,
				(position-currentStopOffset)/float64(nextStopOffset-currentStopOffset),
			)
		}
	}

	return result
}
