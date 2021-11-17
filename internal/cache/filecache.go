package cache

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/atomic"
)

const header = "kitschcache"
const version = uint32(1)

type fileCache struct {
	cacheDir string
}

// TLV types used by the file cache.
const (
	fileCacheValue = 0
	fileCacheKey   = 1
)

// NewFileCache creates a new cache that stores data in the given directory.
// Files are written atomically, so the cache is safe to use across threads and
// processes.  The file cache will hash key names to generate file names, so this
// is safe to use with arbitrary keys.
func NewFileCache(cacheDir string) Cache {
	return &fileCache{cacheDir: cacheDir}
}

// getFileNameForKey returns the full path to the cached object for a given key.
func (cache *fileCache) getFileNameForKey(key string) string {
	// In order to support arbitrary length keys, and keys with special characters
	// in them, we hash the key and use this as the filename.  We take inspiration
	// from git here, and use the first two characters to create a directory.
	hash := sha1.Sum([]byte(key))
	hexHash := hex.EncodeToString(hash[:])
	return filepath.Join(cache.cacheDir, hexHash[0:2], hexHash[2:]+".cache")
}

// Get returns a value from the cache.  Returns nil if the value is not in the
// cache or could not be read.
func (cache *fileCache) Get(key string) []byte {
	fileName := cache.getFileNameForKey(key)
	file, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	defer file.Close()

	cacheKey, cacheValue, err := cache.readCacheRecord(file)
	if err != nil {
		return nil
	}
	if cacheKey != key {
		return nil
	}

	return cacheValue
}

// Set writes a value to the cache.
func (cache *fileCache) Set(key string, value []byte) {
	fileName := cache.getFileNameForKey(key)
	dirName := filepath.Dir(fileName)
	err := os.MkdirAll(dirName, 0700)
	if err != nil {
		return
	}

	file, err := os.CreateTemp(dirName, "temp")
	if err != nil {
		return
	}

	err = cache.writeCacheRecord(file, key, value)
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return
	}

	err = file.Close()
	if err != nil {
		os.Remove(file.Name())
		return
	}

	err = atomic.ReplaceFile(file.Name(), fileName)
	if err != nil {
		os.Remove(file.Name())
	}
}

// Delete removes a value from the cache.
func (cache *fileCache) Delete(key string) {
	fileName := cache.getFileNameForKey(key)
	os.Remove(fileName)
}

func (cache *fileCache) writeCacheRecord(file *os.File, key string, value []byte) error {
	writer := tlvWriter{writer: file}
	writer.WriteHeader([]byte(header), version)
	writer.WriteTLV(fileCacheKey, []byte(key))
	writer.WriteTLV(fileCacheValue, value)
	return writer.End()
}

func (cache *fileCache) readCacheRecord(file *os.File) (key string, value []byte, err error) {
	reader := tlvReader{reader: file}

	recordVersion, err := reader.ReadHeader([]byte(header))
	if err != nil {
		return "", nil, err
	}
	if recordVersion != version {
		return "", nil, fmt.Errorf("incorrect cache version (expected: %d got: %d)", version, recordVersion)
	}

	keyBytes, err := reader.ReadTLVOfType(fileCacheKey)
	if err != nil {
		return "", nil, err
	}
	key = string(keyBytes)

	value, err = reader.ReadTLVOfType(fileCacheValue)
	if err != nil {
		return "", nil, err
	}

	return key, value, nil
}
