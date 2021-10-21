package ansigradient

import (
	"testing"
)

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

func BenchmarkApplyGradients(b *testing.B) {
	var gradient Gradient = CSSLinearGradientMust("#000, #fff")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ApplyGradientsRaw("Hello world!", gradient, nil, LevelAnsi16m)
	}
}
