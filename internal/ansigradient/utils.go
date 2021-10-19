package ansigradient

import (
	"image/color"
)

func linearInterpolateFloat64(start float64, end float64, length int) []float64 {
	answer := make([]float64, length)

	delta := (end - start) / float64(length-1)
	for i := 0; i < length; i++ {
		answer[i] = start + float64(i)*delta
	}

	return answer
}

func lerp(c1 uint8, c2 uint8, s float64) uint8 {
	return uint8(float64(c1) + s*(float64(c2)-float64(c1)))
}

func lerpColor(
	startColor color.RGBA,
	endColor color.RGBA,
	s float64,
) color.RGBA {
	return color.RGBA{
		lerp(startColor.R, endColor.R, s),
		lerp(startColor.G, endColor.G, s),
		lerp(startColor.B, endColor.B, s),
		lerp(startColor.A, endColor.A, s),
	}
}
