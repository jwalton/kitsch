package projects

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/jwalton/kitsch/internal/cache"
	"github.com/jwalton/kitsch/internal/fileutils"
	"github.com/stretchr/testify/assert"
)

type testContext struct {
	fsys  fs.FS
	cache cache.Cache
}

func (t testContext) GetWorkingDirectory() fileutils.Directory {
	return fileutils.NewDirectoryTestFS("/home/jwalton/projects/test", t.fsys)
}

func (testContext) Getenv(key string) string {
	return ""
}

// GetValueCache returns the value cache.
func (t testContext) GetValueCache() cache.Cache {
	if t.cache == nil {
		t.cache = cache.NewMemoryCache()
	}
	return t.cache
}

func TestGetNodeVersionFromVoltaPacakgeJSON(t *testing.T) {
	fsys := fstest.MapFS{
		"package.json": &fstest.MapFile{
			Data: []byte(`{"volta": {"node": "16.0.0"}}`),
		},
	}
	ctx := testContext{fsys: fsys}

	getter := nodejsGetter{executable: "node"}
	version, err := getter.GetValue(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "16.0.0", version)
}
