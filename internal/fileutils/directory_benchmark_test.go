package fileutils

import (
	"testing"
)

func BenchmarkDirectoryHasExtension(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dir := NewDirectory(".")
		dir.HasExtension("go")
	}
}

func BenchmarkDirectoryHasExtensionMultipleCalls(b *testing.B) {
	dir := NewDirectory(".")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dir.HasExtension("go")
	}
}

func BenchmarkDirectoryHasFile(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Include construction of the Directory instance as part of the test,
		// because usually we're not going to check a directory for thousands
		// of files.
		dir := NewDirectory(".")
		dir.HasFile("foo.go")
	}
}

func BenchmarkDirectoryHasGlob(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dir := NewDirectory(".")
		dir.HasGlob("foo.go")
	}
}
