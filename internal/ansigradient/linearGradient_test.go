package ansigradient

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinearGradientColors(t *testing.T) {
	assert.Equal(t,
		[]color.RGBA{
			{R: 14, G: 14, B: 14, A: 255},
			{R: 42, G: 42, B: 42, A: 255},
			{R: 70, G: 70, B: 70, A: 255},
			{R: 99, G: 99, B: 99, A: 255},
			{R: 127, G: 127, B: 127, A: 255},
			{R: 155, G: 155, B: 155, A: 255},
			{R: 184, G: 184, B: 184, A: 255},
			{R: 212, G: 212, B: 212, A: 255},
			{R: 240, G: 240, B: 240, A: 255},
		},
		CSSLinearGradientMust("#000, #fff").Colors(9),
		"Should work if we don't specify stops",
	)

	// This should be *almost* identical to the above.
	assert.Equal(t,
		[]color.RGBA{
			{R: 14, G: 14, B: 14, A: 255},
			{R: 42, G: 42, B: 42, A: 255},
			{R: 70, G: 70, B: 70, A: 255},
			{R: 98, G: 98, B: 98, A: 255},
			{R: 127, G: 127, B: 127, A: 255},
			{R: 155, G: 155, B: 155, A: 255},
			{R: 183, G: 183, B: 183, A: 255},
			{R: 212, G: 212, B: 212, A: 255},
			{R: 240, G: 240, B: 240, A: 255},
		},
		CSSLinearGradientMust("#000, 50%, #fff").Colors(9),
		"Should work if we specify intermediate stop",
	)

	assert.Equal(t,
		[]color.RGBA{
			{R: 14, G: 14, B: 14, A: 255},
			{R: 42, G: 42, B: 42, A: 255},
			{R: 70, G: 70, B: 70, A: 255},
			{R: 99, G: 99, B: 99, A: 255},
			{R: 127, G: 127, B: 127, A: 255},
			{R: 155, G: 155, B: 155, A: 255},
			{R: 184, G: 184, B: 184, A: 255},
			{R: 212, G: 212, B: 212, A: 255},
			{R: 240, G: 240, B: 240, A: 255},
		},
		CSSLinearGradientMust("#000 0%, #fff 100%").Colors(9),
		"Simple gradient",
	)

	assert.Equal(t,
		[]color.RGBA{
			{R: 240, G: 240, B: 240, A: 255},
			{R: 212, G: 212, B: 212, A: 255},
			{R: 184, G: 184, B: 184, A: 255},
			{R: 155, G: 155, B: 155, A: 255},
			{R: 127, G: 127, B: 127, A: 255},
			{R: 99, G: 99, B: 99, A: 255},
			{R: 70, G: 70, B: 70, A: 255},
			{R: 42, G: 42, B: 42, A: 255},
			{R: 14, G: 14, B: 14, A: 255},
		},
		CSSLinearGradientMust("#fff 0%, #000 100%").Colors(9),
		"Reverse gradient",
	)
}

func TestLinearGradientStopsNotAtEnds(t *testing.T) {
	// Should work if stops aren't at start and end
	assert.Equal(t,
		[]color.RGBA{
			{R: 0, G: 0, B: 0, A: 255},
			{R: 0, G: 0, B: 0, A: 255},
			{R: 127, G: 127, B: 127, A: 255},
			{R: 255, G: 255, B: 255, A: 255},
			{R: 255, G: 255, B: 255, A: 255},
		},
		CSSLinearGradientMust("#000 2px, #fff 3px").Colors(5),
		"Should work if stops aren't at start and end",
	)
}
func TestLinearGradientMissingOffsets(t *testing.T) {
	assert.Equal(t,
		[]color.RGBA{
			{R: 0x33, G: 0x33, B: 0x33, A: 0xff},
			{R: 0x99, G: 0x99, B: 0x99, A: 0xff},
			{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			{R: 0x99, G: 0x99, B: 0x99, A: 0xff},
			{R: 0x33, G: 0x33, B: 0x33, A: 0xff},
		},
		CSSLinearGradientMust("#000 0%, #fff, #000 100%").Colors(5),
		"Should interpolate missing offsets",
	)
}

func TestLinearGradientExtraOffsets(t *testing.T) {
	assert.Equal(t,
		[]color.RGBA{
			{R: 0, G: 0, B: 0, A: 255},
			{R: 0, G: 0, B: 0, A: 255},
			{R: 255, G: 255, B: 255, A: 255},
			{R: 0, G: 0, B: 0, A: 255},
			{R: 0, G: 0, B: 0, A: 255},
		},
		CSSLinearGradientMust("#000 0% 2px, #fff, #000 3px 100%").Colors(5),
		"Should correctly render gradient with multiple offsets",
	)

	assert.Equal(t,
		[]color.RGBA{
			{R: 255, G: 0, B: 0, A: 255},
			{R: 255, G: 0, B: 0, A: 255},
			{R: 255, G: 0, B: 0, A: 255},
			{R: 0, G: 0, B: 255, A: 255},
			{R: 0, G: 0, B: 255, A: 255},
			{R: 0, G: 0, B: 255, A: 255},
		},
		CSSLinearGradientMust("#f00 0% 50%, #00f 50% 100%").Colors(6),
		"Should correctly render where first offset is not at 0",
	)
}

func TestLinearGradientAbsoluteStops(t *testing.T) {
	grad := CSSLinearGradientMust("#f00 0px, #00f 5px")

	assert.Equal(
		t,
		[]color.RGBA{
			{R: 0xe5, G: 0, B: 0x19, A: 255},
			{R: 0xb2, G: 0, B: 0x4c, A: 255},
			{R: 0x7f, G: 0, B: 0x7f, A: 255},
			{R: 0x4c, G: 0, B: 0xb2, A: 255},
			{R: 0x19, G: 0, B: 0xe5, A: 255},
			{R: 0x00, G: 0, B: 0xff, A: 255},
		},
		grad.Colors(6),
		"Simple gradient",
	)

	grad = CSSLinearGradientMust("#f00 1px, #00f 5px")
	assert.Equal(t,
		[]color.RGBA{
			{R: 255, G: 0, B: 0, A: 255},
			{R: 223, G: 0, B: 31, A: 255},
			{R: 159, G: 0, B: 95, A: 255},
			{R: 95, G: 0, B: 159, A: 255},
			{R: 31, G: 0, B: 223, A: 255},
		},
		grad.Colors(5),
		"Simple gradient from 1px",
	)

	assert.Equal(t,
		[]color.RGBA{
			{R: 255, G: 0, B: 0, A: 255},
			{R: 223, G: 0, B: 31, A: 255},
			{R: 159, G: 0, B: 95, A: 255},
		},
		grad.Colors(3),
		"Truncated gradient",
	)

	grad = CSSLinearGradientMust("#f00 1px, #00f 1px")
	assert.Equal(t,
		[]color.RGBA{
			{R: 255, G: 0, B: 0, A: 255},
			{R: 0, G: 0, B: 255, A: 255},
			{R: 0, G: 0, B: 255, A: 255},
		},
		grad.Colors(3),
		"Gradient with transition between pixels 1 and 2",
	)
}

func TestLinearGradientOneStop(t *testing.T) {
	grad := CSSLinearGradientMust("#f00")

	assert.Equal(
		t,
		[]color.RGBA{
			{R: 0xff, G: 0, B: 0x0, A: 255},
			{R: 0xff, G: 0, B: 0x0, A: 255},
		},
		grad.Colors(2),
	)
}
