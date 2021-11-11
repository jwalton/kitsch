package getters

import (
	"testing"
	"testing/fstest"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/stretchr/testify/assert"
)

func TestYamlGetterFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"version.txt": &fstest.MapFile{
			Data: []byte("v1.0.0\n"),
		},
	}

	directory := fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	getter := CustomGetter{
		Type: "file",
		From: "version.txt",
	}

	val, err := getter.GetValue(directory)

	assert.Nil(t, err)
	assert.Equal(t, "v1.0.0", val)
}

func TestYamlGetterTextFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"version.txt": &fstest.MapFile{
			Data: []byte("v1.0.0\n"),
		},
	}

	directory := fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	// Verify if we specify "As: text" with no template or regex, we get the
	// right result.
	getter := CustomGetter{
		Type: "file",
		From: "version.txt",
		As:   "text",
	}

	val, err := getter.GetValue(directory)

	assert.Nil(t, err)
	assert.Equal(t, "v1.0.0", val)
}

func TestYamlGetterJsonFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"version.json": &fstest.MapFile{
			Data: []byte(`{"version": "v1.0.0"}` + "\n"),
		},
	}

	directory := fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	getter := CustomGetter{
		Type:          "file",
		From:          "version.json",
		As:            "json",
		ValueTemplate: `{{.version}}`,
	}

	val, err := getter.GetValue(directory)

	assert.Nil(t, err)
	assert.Equal(t, "v1.0.0", val)
}

func TestYamlGetterJsonFromFileNoTemplate(t *testing.T) {
	fsys := fstest.MapFS{
		"version.json": &fstest.MapFile{
			Data: []byte(`{"version": "v1.0.0"}` + "\n"),
		},
	}

	directory := fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	getter := CustomGetter{
		Type: "file",
		From: "version.json",
		As:   "json",
	}

	val, err := getter.GetValue(directory)

	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"version": "v1.0.0"}, val)
}

func TestYamlGetterYamlFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"version.yaml": &fstest.MapFile{
			Data: []byte(`version: v1.0.0`),
		},
	}

	directory := fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	getter := CustomGetter{
		Type:          "file",
		From:          "version.yaml",
		As:            "yaml",
		ValueTemplate: `{{.version}}`,
	}

	val, err := getter.GetValue(directory)

	assert.Nil(t, err)
	assert.Equal(t, "v1.0.0", val)
}

func TestYamlGetterTomlFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"version.toml": &fstest.MapFile{
			Data: []byte(`version = "v1.0.0"`),
		},
	}

	directory := fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	getter := CustomGetter{
		Type:          "file",
		From:          "version.toml",
		As:            "toml",
		ValueTemplate: `{{.version}}`,
	}

	val, err := getter.GetValue(directory)

	assert.Nil(t, err)
	assert.Equal(t, "v1.0.0", val)
}

func TestYamlGetterNumericValue(t *testing.T) {
	fsys := fstest.MapFS{
		"version.yaml": &fstest.MapFile{
			Data: []byte(`version: 10`),
		},
	}

	directory := fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	getter := CustomGetter{
		Type:          "file",
		From:          "version.yaml",
		As:            "yaml",
		ValueTemplate: `{{.version}}`,
	}

	val, err := getter.GetValue(directory)

	assert.Nil(t, err)
	assert.Equal(t, "10", val)
}

func TestYamlGetterRegex(t *testing.T) {
	fsys := fstest.MapFS{
		"version.txt": &fstest.MapFile{
			Data: []byte("go version go1.17.1 darwin/amd64\n"),
		},
	}

	directory := fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	getter := CustomGetter{
		Type:  "file",
		From:  "version.txt",
		Regex: `go version go(\d+\.\d+\.\d+)`,
	}

	val, err := getter.GetValue(directory)

	assert.Nil(t, err)
	assert.Equal(t, "1.17.1", val)
}

func TestYamlGetterRegexAndTemplate(t *testing.T) {
	fsys := fstest.MapFS{
		"version.txt": &fstest.MapFile{
			Data: []byte("go version go1.17.1 darwin/amd64\n"),
		},
	}

	directory := fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	getter := CustomGetter{
		Type:          "file",
		From:          "version.txt",
		Regex:         `go version go(\d+\.\d+\.\d+)`,
		ValueTemplate: `v{{.Text}}`,
	}

	val, err := getter.GetValue(directory)

	assert.Nil(t, err)
	assert.Equal(t, "v1.17.1", val)
}

func TestYamlGetterTextTemplate(t *testing.T) {
	fsys := fstest.MapFS{
		"version.txt": &fstest.MapFile{
			Data: []byte("1.17.1\n"),
		},
	}

	directory := fileutils.NewDirectoryTestFS("/foo/bar", fsys)

	getter := CustomGetter{
		Type:          "file",
		From:          "version.txt",
		ValueTemplate: `v{{.Text}}`,
	}

	val, err := getter.GetValue(directory)

	assert.Nil(t, err)
	assert.Equal(t, "v1.17.1", val)
}
