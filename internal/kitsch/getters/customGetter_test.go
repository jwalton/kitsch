package getters

import (
	"testing"
	"testing/fstest"

	"github.com/jwalton/kitsch-prompt/internal/cache"
	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/stretchr/testify/assert"
)

type testGetterContext struct {
	directory fileutils.Directory
	cache     cache.Cache
}

// GetWorkingDirectory returns the current working directory.
func (context *testGetterContext) GetWorkingDirectory() fileutils.Directory {
	return context.directory
}

// Getenv returns the value of the specified environment variable.
func (context *testGetterContext) Getenv(key string) string {
	return ""
}

// GetValueCache returns the value cache.
func (context *testGetterContext) GetValueCache() cache.Cache {
	return context.cache
}

func makeTestGetterContext(fsys fstest.MapFS) *testGetterContext {
	return &testGetterContext{
		directory: fileutils.NewDirectoryTestFS("/foo/bar", fsys),
		cache:     cache.NewMemoryCache(),
	}
}

func TestYamlGetterFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"version.txt": &fstest.MapFile{
			Data: []byte("v1.0.0\n"),
		},
	}
	context := makeTestGetterContext(fsys)

	getter := CustomGetter{
		Type: TypeFile,
		From: "version.txt",
	}

	val, err := getter.GetValue(context)

	assert.Nil(t, err)
	assert.Equal(t, "v1.0.0", val)
}

func TestYamlGetterTextFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"version.txt": &fstest.MapFile{
			Data: []byte("v1.0.0\n"),
		},
	}
	context := makeTestGetterContext(fsys)

	// Verify if we specify "As: text" with no template or regex, we get the
	// right result.
	getter := CustomGetter{
		Type: TypeFile,
		From: "version.txt",
		As:   AsText,
	}

	val, err := getter.GetValue(context)

	assert.Nil(t, err)
	assert.Equal(t, "v1.0.0", val)
}

func TestYamlGetterJsonFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"version.json": &fstest.MapFile{
			Data: []byte(`{"version": "v1.0.0"}` + "\n"),
		},
	}
	context := makeTestGetterContext(fsys)

	getter := CustomGetter{
		Type:          TypeFile,
		From:          "version.json",
		As:            AsJSON,
		ValueTemplate: `{{.version}}`,
	}

	val, err := getter.GetValue(context)

	assert.Nil(t, err)
	assert.Equal(t, "v1.0.0", val)
}

func TestYamlGetterJsonFromFileNoTemplate(t *testing.T) {
	fsys := fstest.MapFS{
		"version.json": &fstest.MapFile{
			Data: []byte(`{"version": "v1.0.0"}` + "\n"),
		},
	}
	context := makeTestGetterContext(fsys)

	getter := CustomGetter{
		Type: TypeFile,
		From: "version.json",
		As:   AsJSON,
	}

	val, err := getter.GetValue(context)

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"version": "v1.0.0"}, val)
}

func TestYamlGetterYamlFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"version.yaml": &fstest.MapFile{
			Data: []byte(`version: v1.0.0`),
		},
	}
	context := makeTestGetterContext(fsys)

	getter := CustomGetter{
		Type:          TypeFile,
		From:          "version.yaml",
		As:            AsYAML,
		ValueTemplate: `{{.version}}`,
	}

	val, err := getter.GetValue(context)

	assert.Nil(t, err)
	assert.Equal(t, "v1.0.0", val)
}

func TestYamlGetterTomlFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"version.toml": &fstest.MapFile{
			Data: []byte(`version = "v1.0.0"`),
		},
	}
	context := makeTestGetterContext(fsys)

	getter := CustomGetter{
		Type:          TypeFile,
		From:          "version.toml",
		As:            AsTOML,
		ValueTemplate: `{{.version}}`,
	}

	val, err := getter.GetValue(context)

	assert.Nil(t, err)
	assert.Equal(t, "v1.0.0", val)
}

func TestYamlGetterNumericValue(t *testing.T) {
	fsys := fstest.MapFS{
		"version.yaml": &fstest.MapFile{
			Data: []byte(`version: 10`),
		},
	}
	context := makeTestGetterContext(fsys)

	getter := CustomGetter{
		Type:          TypeFile,
		From:          "version.yaml",
		As:            AsYAML,
		ValueTemplate: `{{.version}}`,
	}

	val, err := getter.GetValue(context)

	assert.Nil(t, err)
	assert.Equal(t, "10", val)
}

func TestYamlGetterRegex(t *testing.T) {
	fsys := fstest.MapFS{
		"version.txt": &fstest.MapFile{
			Data: []byte("go version go1.17.1 darwin/amd64\n"),
		},
	}
	context := makeTestGetterContext(fsys)

	getter := CustomGetter{
		Type:  TypeFile,
		From:  "version.txt",
		Regex: `go version go(\d+\.\d+\.\d+)`,
	}

	val, err := getter.GetValue(context)

	assert.Nil(t, err)
	assert.Equal(t, "1.17.1", val)
}

func TestYamlGetterRegexAndTemplate(t *testing.T) {
	fsys := fstest.MapFS{
		"version.txt": &fstest.MapFile{
			Data: []byte("go version go1.17.1 darwin/amd64\n"),
		},
	}
	context := makeTestGetterContext(fsys)

	getter := CustomGetter{
		Type:          TypeFile,
		From:          "version.txt",
		Regex:         `go version go(\d+\.\d+\.\d+)`,
		ValueTemplate: `v{{.Text}}`,
	}

	val, err := getter.GetValue(context)

	assert.Nil(t, err)
	assert.Equal(t, "v1.17.1", val)
}

func TestYamlGetterTextTemplate(t *testing.T) {
	fsys := fstest.MapFS{
		"version.txt": &fstest.MapFile{
			Data: []byte("1.17.1\n"),
		},
	}
	context := makeTestGetterContext(fsys)

	getter := CustomGetter{
		Type:          TypeFile,
		From:          "version.txt",
		ValueTemplate: `v{{.Text}}`,
	}

	val, err := getter.GetValue(context)

	assert.Nil(t, err)
	assert.Equal(t, "v1.17.1", val)
}
