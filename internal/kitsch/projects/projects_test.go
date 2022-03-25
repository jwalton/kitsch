package projects

import (
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/jwalton/kitsch/internal/cache"
	"github.com/jwalton/kitsch/internal/fileutils"
	"github.com/jwalton/kitsch/internal/kitsch/getters"
	"github.com/stretchr/testify/assert"
)

type testGetter struct {
	value     string
	shouldErr bool
}

func (g *testGetter) GetValue(getterContext getters.GetterContext) (interface{}, error) {
	if g.shouldErr {
		return nil, fmt.Errorf("error")
	}
	return g.value, nil
}

type testGetterContext struct {
	directory fileutils.Directory
	home      string
	cache     cache.Cache
	env       map[string]string
}

// GetWorkingDirectory returns the current working directory.
func (context *testGetterContext) GetWorkingDirectory() fileutils.Directory {
	return context.directory
}

// GetHomeDirectoryPath returns the path to the user's home directory.
func (context *testGetterContext) GetHomeDirectoryPath() string {
	return context.home
}

// Getenv returns the value of the specified environment variable.
func (context *testGetterContext) Getenv(key string) string {
	return context.env[key]
}

// GetValueCache returns the value cache.
func (context *testGetterContext) GetValueCache() cache.Cache {
	return context.cache
}

func makeTestGetterContext(fsys fstest.MapFS) *testGetterContext {
	return &testGetterContext{
		directory: fileutils.NewDirectoryTestFS("/foo/bar", fsys),
		home:      "/users/jwalton",
		cache:     cache.NewMemoryCache(),
	}
}

func TestGetStringValueNoGetters(t *testing.T) {
	result, err := getStringValue([]getters.Getter{}, makeTestGetterContext(fstest.MapFS{}))

	assert.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestGetStringValue(t *testing.T) {
	result, err := getStringValue([]getters.Getter{
		&testGetter{value: "foo"},
	}, makeTestGetterContext(fstest.MapFS{}))

	assert.NoError(t, err)
	assert.Equal(t, "foo", result)
}

func TestGetStringErr(t *testing.T) {
	_, err := getStringValue([]getters.Getter{
		&testGetter{shouldErr: true},
	}, makeTestGetterContext(fstest.MapFS{}))

	assert.EqualError(t, err, "error")
}

func TestGetStringMultiple(t *testing.T) {
	result, err := getStringValue([]getters.Getter{
		&testGetter{shouldErr: true},
		&testGetter{value: "foo"},
	}, makeTestGetterContext(fstest.MapFS{}))

	assert.NoError(t, err)
	assert.Equal(t, "foo", result)
}
