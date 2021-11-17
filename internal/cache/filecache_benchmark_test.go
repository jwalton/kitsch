package cache

import (
	"os"
	"testing"
)

func BenchmarkFileCacheNames(b *testing.B) {
	tempdir, err := os.MkdirTemp("", "filecache")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	cache := fileCache{cacheDir: tempdir}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.getFileNameForKey("/Users/jwalton/.nvm/versions/node/v14.18.1/lib/node_modules/npm/bin/npm-cli.js")
	}
}
