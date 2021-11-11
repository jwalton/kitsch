package projects

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
)

type npmVersionGetter struct{}

// GetValue returns the version of npm.
//
// Running `npm --version` is quite slow, so we cheat a little here and find
// the `npm` executable, then go find the associated `package.json` and read
// the version from there.
func (getter npmVersionGetter) GetValue(folder fileutils.Directory) (interface{}, error) {
	npmPath, err := fileutils.LookPathSafe("npm")
	if err != nil {
		return "", err
	}

	var version string
	if runtime.GOOS == "windows" {
		version, err = getter.getNpmVersionWindows(npmPath)
	} else {
		version, err = getter.getNpmVersion(npmPath)
	}

	if err != nil && npmPath != "" {
		// If the hacky version failed, then just run `npm --version`
		cmd := exec.Command(npmPath, "--version")
		cmd.Dir = folder.Path()

		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}

		version = string(out)
	}

	return version, err
}

func (getter npmVersionGetter) getNpmVersion(npmPath string) (string, error) {
	// Find the "npm" command - it will be a symlink.  Follow this to the npm executable, then
	// walk up the directory tree to find the package.json.
	npmPath, err := filepath.EvalSymlinks(npmPath)
	if err != nil {
		return "", err
	}

	npmDir := filepath.Dir(npmPath)
	if strings.HasSuffix(npmDir, "node_modules/npm/bin") {
		packageJSONPath := filepath.Join(filepath.Dir(npmDir), "package.json")
		return getter.readVersionFromPackageJSON(packageJSONPath)
	}

	return "", fmt.Errorf("could not find npm version")
}

func (getter npmVersionGetter) getNpmVersionWindows(npmPath string) (string, error) {
	packageJSONPath := filepath.Join(filepath.Dir(npmPath), "node_modules", "npm", "package.json")
	return getter.readVersionFromPackageJSON(packageJSONPath)
}

func (getter npmVersionGetter) readVersionFromPackageJSON(packageJSONPath string) (string, error) {
	bytes, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return "", err
	}

	var packageJSON map[string]interface{}
	err = json.Unmarshal(bytes, &packageJSON)
	if err != nil {
		return "", err
	}

	version, ok := packageJSON["version"]
	if !ok {
		return "", fmt.Errorf("package.json does not contain a version")
	}

	return version.(string), nil
}
