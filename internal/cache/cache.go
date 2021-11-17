// Package cache implements a simple file-based cache, which can be accessed
// from multiple concurrent processes.  Files are cached in a cache directory,
// in a file based on the key name.  When a value is written to the cache, it is
// written "atomically" by writing it to a file then renaming the file to the
// final destination.  If a file is corrupt when reading, we return a cache miss.
package cache

// Cache represents a cache for values.
type Cache interface {
	// Get returns the value for the given key.  If the value is not found,
	// returns nil.
	Get(key string) []byte
	// Set sets the value for the given key.
	Set(key string, value []byte)
	// Delete deletes the value for the given key.
	Delete(key string)
}
