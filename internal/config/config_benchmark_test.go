package config

import (
	"testing"
)

func BenchmarkLoadDefaultConfig(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = LoadDefaultConfig()
	}
}
