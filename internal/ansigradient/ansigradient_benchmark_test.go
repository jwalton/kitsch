package ansigradient

import (
	"testing"
)

func BenchmarkColorStringFromGradient(b *testing.B) {
	var gradient Gradient = CSSLinearGradientMust("#000, #fff")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ApplyGradientsRaw("Hello world!", gradient, nil, LevelAnsi16m)
	}
}

func BenchmarkGradient(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		CSSLinearGradientMust("#000, #fff")
	}
}

func BenchmarkGradientGenerateColors(b *testing.B) {
	var gradient Gradient = CSSLinearGradientMust("#000, #fff")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		gradient.Colors(10)
	}
}

func BenchmarkColorString(b *testing.B) {
	var gradient Gradient = CSSLinearGradientMust("#000, #fff")
	colors := gradient.Colors(10)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ColorStringRaw("Hello world!", colors, nil, LevelAnsi16m)
	}
}
