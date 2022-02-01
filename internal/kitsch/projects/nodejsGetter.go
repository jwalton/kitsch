package projects

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jwalton/kitsch/internal/fileutils"
	"github.com/jwalton/kitsch/internal/kitsch/getters"
)

// nodejsGetter is a volta-aware getter for the current node/npm version.
//
// When we run `node --version`, we usually let the CustomGetter try to find the
// node executable and cache the the path, size, and timestamp of the node executable.
// If the user is using volta to manage their node/npm version, though, this doesn't
// work, because the "node" symlink goes to ~/.volta/bin/volta-shim, which is
// the same executable no matter which version of node we're using.  One solution
// here would be to disable caching for node.js, but `npm --version` is crazy slow
// to run - about half a second - so we really don't want to disable the cache here.
//
// Instead we have this custom "nodejsGetter" which tries to detect if we're using
// volta, and if so we run `volta which node` or `volta which npm` to work out what
// version we're running.
type nodejsGetter struct {
	// executable should be "node", "npm", or "yarn".
	executable string
	regex      string
}

func (getter nodejsGetter) GetValue(getterContext getters.GetterContext) (interface{}, error) {
	// Resolve the executable to an absolute path.
	executable, err := fileutils.LookPathSafe(getter.executable)
	if err != nil {
		return nil, fmt.Errorf("could not find executable: \"%s\": %w", getter.executable, err)
	}

	// If the executable is a symlink, resolve it.
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return nil, fmt.Errorf("could not resolve executable: \"%s\": %w", getter.executable, err)
	}

	// If the executable is "volta-shim", ask volta for the command we're actually going to run.
	if strings.HasSuffix(executable, "volta-shim") || strings.HasSuffix(executable, "volta-shim.exe") {
		return getter.getFromVolta(getterContext)
	}

	return getter.getFromExecutable(getterContext, executable)
}

func (getter nodejsGetter) getFromVolta(getterContext getters.GetterContext) (interface{}, error) {
	result, err := getter.getFromVoltaPackageJSON(getterContext)
	if err == nil {
		return result, nil
	}

	result, err = getter.getFromVoltaConfig(getterContext)
	if err == nil {
		return result, nil
	}

	// If the version isn't in package.json, and we can't get it from the Volta config, then we need to run volta.
	volta, err := fileutils.LookPathSafe("volta")
	if err != nil {
		return nil, fmt.Errorf("could not find volta: %w", err)
	}

	cmd := exec.Command(volta, "which", getter.executable)
	cmd.Dir = getterContext.GetWorkingDirectory().Path()
	executable, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("could not resolve \"%s\" target version: %w", getter.executable, err)
	}

	return getter.getFromExecutable(getterContext, strings.TrimSpace(string(executable)))
}

func (getter nodejsGetter) getFromVoltaPackageJSON(getterContext getters.GetterContext) (string, error) {
	// First, try to parse the version from the volta section of package.json.
	rawPackageJSON, err := fs.ReadFile(getterContext.GetWorkingDirectory().FileSystem(), "package.json")
	if err != nil {
		return "", fmt.Errorf("could not read package.json: %w", err)
	}

	packageJSON := map[string]interface{}{}
	err = json.Unmarshal(rawPackageJSON, &packageJSON)
	if err != nil {
		return "", fmt.Errorf("could not parse package.json: %w", err)
	}

	voltaSection, ok := packageJSON["volta"]
	if !ok {
		return "", fmt.Errorf("package.json does not have a volta section")
	}
	voltaSectionMap, ok := voltaSection.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("volta section in package.json is not expected type")
	}
	version, ok := voltaSectionMap[getter.executable]
	if !ok {
		return "", fmt.Errorf("%s missing in volta section", getter.executable)
	}
	result, ok := version.(string)
	if !ok {
		return "", fmt.Errorf("%s is in volta section but is not a string", getter.executable)
	}

	return result, nil
}

type voltaConfig struct {
	Node struct {
		Runtime string `json:"runtime"`
		Npm     string `json:"npm"`
	} `json:"node"`
	Yarn string `json:"yarn"`
}

// Try to get the nodejs version from ~/.volta/tools/user/platform.json
func (getter nodejsGetter) getFromVoltaConfig(getterContext getters.GetterContext) (string, error) {
	voltaHome := getterContext.Getenv("VOLTA_HOME")
	if voltaHome == "" {
		voltaHome = filepath.Join(getterContext.Getenv("HOME") + ".volta")
	}
	configPath := filepath.Join(voltaHome, "tools", "user", "platform.json")
	configContents, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("could not read %s: %w", configPath, err)
	}

	var config voltaConfig
	err = json.Unmarshal(configContents, &config)
	if err != nil {
		return "", fmt.Errorf("could not parse %s: %w", configPath, err)
	}

	switch getter.executable {
	case "node":
		if config.Node.Runtime == "" {
			return "", fmt.Errorf("No node version specified")
		}
		return config.Node.Runtime, nil
	case "npm":
		if config.Node.Npm == "" {
			return "", fmt.Errorf("No npm version specified")
		}
		return config.Node.Npm, nil
	case "yarn":
		if config.Yarn == "" {
			return "", fmt.Errorf("No yarn version specified")
		}
		return config.Yarn, nil
	default:
		return "", fmt.Errorf("Unknown executable: %s", getter.executable)
	}
}

func (getter nodejsGetter) getFromExecutable(getterContext getters.GetterContext, resolvedExecutable string) (interface{}, error) {
	// Delegate this to a custom getter to handle caching.
	customGetter := getters.CustomGetter{
		Type:  getters.TypeCustom,
		From:  resolvedExecutable + " --version",
		Regex: getter.regex,
		Cache: getters.CacheSettings{Enabled: true},
	}

	return customGetter.GetValue(getterContext)
}
