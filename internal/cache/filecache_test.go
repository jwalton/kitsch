package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileCache(t *testing.T) {
	tempdir, err := os.MkdirTemp("", "filecache")
	assert.NoError(t, err)
	defer os.RemoveAll(tempdir)

	cache := NewFileCache(tempdir)

	// Write a value to the cache.
	cache.Set("foo", []byte("bar"))

	// Make sure we can read the value back.
	result := cache.Get("foo")
	assert.Equal(t, []byte("bar"), result)

	// Make sure getting a value that isn't in the cache returns nil.
	result = cache.Get("i.do.not.exist")
	assert.Nil(t, result)

	// Delete the key
	cache.Delete("foo")

	// Make sure it's gone.
	result = cache.Get("foo")
	assert.Nil(t, result)
}

func TestFileColission(t *testing.T) {
	tempdir, err := os.MkdirTemp("", "filecache")
	assert.NoError(t, err)
	defer os.RemoveAll(tempdir)

	cache := fileCache{cacheDir: tempdir}

	key := "foo"
	filename := cache.getFileNameForKey(key)

	// Write a file with a different key at this location.
	err = os.MkdirAll(filepath.Dir(filename), 0700)
	assert.NoError(t, err)

	file, err := os.Create(filename)
	assert.NoError(t, err)
	defer file.Close()

	err = cache.writeCacheRecord(file, "qux", []byte("bar"))
	assert.NoError(t, err)

	// Try to get the value back.  Should fail because key doesn't match.
	result := cache.Get(key)
	assert.Nil(t, result)
}

func TestFileCacheNames(t *testing.T) {
	tempdir, err := os.MkdirTemp("", "filecache")
	assert.NoError(t, err)
	defer os.RemoveAll(tempdir)

	cache := fileCache{cacheDir: tempdir}

	filename := cache.getFileNameForKey("foo")
	assert.Equal(t, filepath.Join(tempdir, "0b", "eec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33.cache"), filename)

	filename = cache.getFileNameForKey("foo/bar")
	assert.Equal(t, filepath.Join(tempdir, "17", "cdeaefa5cc6022481c824e15a47a7726f593dd.cache"), filename)
}
