package fileutils

import (
	"testing"
)

func BenchmarkFindFileInAncestors(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FindFileInAncestors(".", "foo.txt")
	}
}

func BenchmarkFileExists(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fileExists("foo.txt")
	}
}
