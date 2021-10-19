package ansigradient

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCssStopParser_SingleStops(t *testing.T) {
	result, err := parseCSSStops(nil, "#010203")
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]gradientStop{{Color: color.RGBA{1, 2, 3, 255}, ColorUnset: false, Offset: 0, OffsetType: gradientStopUnspecified}},
		result,
	)

	result, err = parseCSSStops(nil, "#010203 10%")
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]gradientStop{{Color: color.RGBA{1, 2, 3, 255}, ColorUnset: false, Offset: 0.1, OffsetType: gradientStopRelative}},
		result,
	)

	result, err = parseCSSStops(nil, "#010203 10px")
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]gradientStop{{Color: color.RGBA{1, 2, 3, 255}, ColorUnset: false, Offset: 10, OffsetType: gradientStopAbsolute}},
		result,
	)

	result, err = parseCSSStops(nil, "#010203 10px 20px")
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]gradientStop{
			{Color: color.RGBA{1, 2, 3, 255}, ColorUnset: false, Offset: 10, OffsetType: gradientStopAbsolute},
			{Color: color.RGBA{1, 2, 3, 255}, ColorUnset: false, Offset: 20, OffsetType: gradientStopAbsolute},
		},
		result,
	)

	result, err = parseCSSStops(nil, "#010203 10% 20px")
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]gradientStop{
			{Color: color.RGBA{1, 2, 3, 255}, ColorUnset: false, Offset: 0.1, OffsetType: gradientStopRelative},
			{Color: color.RGBA{1, 2, 3, 255}, ColorUnset: false, Offset: 20, OffsetType: gradientStopAbsolute},
		},
		result,
	)

	result, err = parseCSSStops(nil, "#010203 10px 20%")
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]gradientStop{
			{Color: color.RGBA{1, 2, 3, 255}, ColorUnset: false, Offset: 10, OffsetType: gradientStopAbsolute},
			{Color: color.RGBA{1, 2, 3, 255}, ColorUnset: false, Offset: 0.2, OffsetType: gradientStopRelative},
		},
		result,
	)

	// Should accept "0" with no % or px
	result, err = parseCSSStops(nil, "#010203 0")
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]gradientStop{{Color: color.RGBA{1, 2, 3, 255}, ColorUnset: false, Offset: 0, OffsetType: gradientStopAbsolute}},
		result,
	)
}

func TestCssStopParser_MultiStops(t *testing.T) {
	result, err := parseCSSStops(nil, "#fff, #000, #fff")
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]gradientStop{
			{Color: color.RGBA{255, 255, 255, 255}, ColorUnset: false, Offset: 0, OffsetType: gradientStopUnspecified},
			{Color: color.RGBA{0, 0, 0, 255}, ColorUnset: false, Offset: 0, OffsetType: gradientStopUnspecified},
			{Color: color.RGBA{255, 255, 255, 255}, ColorUnset: false, Offset: 0, OffsetType: gradientStopUnspecified},
		},
		result,
	)

	result, err = parseCSSStops(nil, "#f00 10px 20%, 30%, #00f")
	assert.Nil(t, err)
	assert.Equal(
		t,
		[]gradientStop{
			{Color: color.RGBA{255, 0, 0, 255}, ColorUnset: false, Offset: 10, OffsetType: gradientStopAbsolute},
			{Color: color.RGBA{255, 0, 0, 255}, ColorUnset: false, Offset: 0.2, OffsetType: gradientStopRelative},
			{Color: color.RGBA{}, ColorUnset: true, Offset: 0.3, OffsetType: gradientStopRelative},
			{Color: color.RGBA{0, 0, 255, 255}, ColorUnset: false, Offset: 0, OffsetType: gradientStopUnspecified},
		},
		result,
	)
}

func TestCssStopParser_BadStops(t *testing.T) {
	_, err := parseCSSStops(nil, "10%, #fff")
	assert.Equal(
		t,
		"Expected linear-color-stop at position 1",
		err.Error(),
	)

	_, err = parseCSSStops(nil, "#fff, 10%, 90%, #fff")
	assert.Equal(
		t,
		"Cannot have two linear-color-hint in a row at position 12",
		err.Error(),
	)

	_, err = parseCSSStops(nil, "#22, #fff")
	assert.Equal(
		t,
		"invalid color at 0: invalid hex color \"#22\"",
		err.Error(),
	)

	_, err = parseCSSStops(nil, "#222, woo")
	assert.Equal(
		t,
		"expected color at position 7, got \"woo\"",
		err.Error(),
	)

	_, err = parseCSSStops(nil, "#222 #fff")
	assert.Equal(
		t,
		"expected ',' at position 6, found '#'",
		err.Error(),
	)
}
