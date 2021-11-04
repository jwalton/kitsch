package fileutils

import (
	"testing"
)

func BenchmarkFindFileInAncestors(b *testing.B) {
	utils := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utils.FindFileInAncestors(".", "foo.txt")
	}
}

func BenchmarkFileExists(b *testing.B) {
	utils := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utils.FileExists("foo.txt")
	}
}
