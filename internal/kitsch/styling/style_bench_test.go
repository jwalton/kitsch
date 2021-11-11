package styling

import (
	"testing"
)

func BenchmarkParseStyle(b *testing.B) {
	styles := testStyleRegistry()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := styles.Get("red")
		if err != nil {
			b.Fatal(err)
		}

		// Clear the style from the cache.
		delete(styles.styles, "red")
	}
}

func BenchmarkApplyStyle(b *testing.B) {
	styles := testStyleRegistry()
	style, err := styles.Get("red")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		style.Apply("test")
	}
}
