package getters

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig/v3"
	"github.com/jwalton/kitsch-prompt/internal/cache"
	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/modtemplate"
	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v3"
)

var sprigTemplateFunctions = sprig.TxtFuncMap()

//go:generate stringer -type=GetterType

// GetterType is the type of getter.
type GetterType int

const (
	// TypeCustom is a getter which runs a command and returns the output.
	TypeCustom GetterType = iota
	// TypeFile is a getter which reads a file and returns the contents.
	TypeFile
	// TypeAncestorFile is a getter which reads a file in the current folder or
	// any ancestor folder and returns the contents.
	TypeAncestorFile
	// TypeEnv is a getter which reads an environment variable and returns the value.
	TypeEnv
)

// UnmarshalYAML unmarshals a YAML node into a GetterType.
func (item *GetterType) UnmarshalYAML(node *yaml.Node) error {
	var value string
	if err := node.Decode(&value); err != nil {
		return err
	}

	switch value {
	case "custom":
		*item = TypeCustom
	case "file":
		*item = TypeFile
	case "ancestorFile":
		*item = TypeAncestorFile
	case "env":
		*item = TypeEnv
	default:
		return fmt.Errorf("unknown GetterType: %s", value)
	}

	return nil
}

//go:generate stringer -type=AsType

// AsType describes how to interpret the retrieved value.
type AsType int

const (
	// AsUndefined will parse the returned value as the default type.
	AsUndefined AsType = iota
	// AsText will parse the returned value as a string.
	AsText
	// AsJSON will parse the returned value as a JSON object.
	AsJSON
	// AsYAML will parse the returned value as a YAML file.
	AsYAML
	// AsTOML will parse the returned value as a TOML file.
	AsTOML
)

// UnmarshalYAML unmarshals a YAML node into an AsType.
func (item *AsType) UnmarshalYAML(node *yaml.Node) error {
	var value string
	if err := node.Decode(&value); err != nil {
		return err
	}

	switch value {
	case "text":
		*item = AsText
	case "json":
		*item = AsJSON
	case "yaml":
		*item = AsYAML
	case "toml":
		*item = AsTOML
	default:
		return fmt.Errorf("unknown AsType: %s", value)
	}

	return nil
}

// CacheSettings are cache settings for a CustomGetter.
type CacheSettings struct {
	// Enabled is true if caching should be enabled for this getter.
	//
	// At the moment, this only applied to getters with `Type: "custom"`.  This
	// makes it so we will cache the output of a command instead of re-running that
	// command.
	Enabled bool `yaml:"enabled"`
}

// CustomGetter is a getter that can be configured from a YAML file.
type CustomGetter struct {
	// Type is the type of getter.  One of "custom", "file", "ancestorFile", or "env".
	Type GetterType `yaml:"type"`
	// From is the source to get data from.  The meaning of "From" is based on
	// the provided "Type".
	From string `yaml:"from"`
	// As will determine how to interpret the result of the getter.  One of "text", "json", "toml", or "yaml".
	As AsType `yaml:"as"`
	// ValueTemplate is a golang template used to parse values out of the result of
	// the getter.
	ValueTemplate string `yaml:"valueTemplate"`
	// Regex is a regular expression used to parse values out of the result of
	// the getter.  If specified, then "As" and "Template" will be ignored.
	Regex string `yaml:"regex"`
	// Cache specified cache settings for this getter.
	Cache CacheSettings `yaml:"cache"`
}

// GetValue gets the value for this getter.  The return value will be either a string,
// of if the value is a JSON, YAML, or TOML object, and the `ValueTemplate` is not set,
// the parsed contents of the object.
func (getter CustomGetter) GetValue(context GetterContext) (interface{}, error) {
	// Get the raw value for the getter.
	var bytesValue []byte
	var err error
	var result interface{}

	folder := context.GetWorkingDirectory()
	valueCache := context.GetValueCache()

	switch getter.Type {
	case TypeCustom:
		bytesValue, err = getter.getCustomValue(folder, valueCache, getter.From)
	case TypeFile:
		bytesValue, err = fs.ReadFile(folder.FileSystem(), getter.From)
	case TypeAncestorFile:
		bytesValue, err = getter.getAncestorFileValue(folder, getter.From)
	case TypeEnv:
		strValue := context.Getenv(getter.From)
		if strValue == "" {
			bytesValue = nil
		} else {
			bytesValue = []byte(strValue)
		}
	default:
		err = fmt.Errorf("invalid getter type: \"%s\"", getter.Type)
	}
	if err != nil {
		return "", err
	}
	if bytesValue == nil {
		return "", nil
	}

	// Run the value through the regex, if required.
	if getter.Regex != "" {
		regex, err := regexp.Compile(getter.Regex)
		if err != nil {
			return "", fmt.Errorf("invalid regex: \"%s\": %w", getter.Regex, err)
		}

		matches := regex.FindStringSubmatch(string(bytesValue))

		var strValue string
		if len(matches) == 0 {
			strValue = ""
		} else if len(matches) > 1 {
			strValue = matches[1]
		} else {
			strValue = matches[0]
		}

		if getter.ValueTemplate != "" {
			result, err = getter.applyTemplate(AsText, []byte(strValue))
			if err != nil {
				return "", err
			}
		} else {
			result = strValue
		}

	} else if (getter.As != AsUndefined || getter.ValueTemplate != "") && !(getter.As == AsText && getter.ValueTemplate == "") {
		as := getter.As
		if as == AsUndefined {
			as = AsText
		}
		result, err = getter.applyTemplate(as, bytesValue)
		if err != nil {
			return "", err
		}

	} else {
		result = strings.TrimSpace(string(bytesValue))
	}

	return result, nil
}

func (getter CustomGetter) getCustomValue(
	projectFolder fileutils.Directory,
	valueCache cache.Cache,
	command string,
) ([]byte, error) {
	commandParts, err := shellwords.Parse(command)
	if err != nil {
		return nil, fmt.Errorf("invalid command: \"%s\": %w", command, err)
	}
	if len(commandParts) == 0 {
		return nil, fmt.Errorf("invalid command: \"%s\"", command)
	}

	// Resolve the executable to an absolute path.
	executable := commandParts[0]
	executable, err = fileutils.LookPathSafe(executable)
	if err != nil {
		return nil, fmt.Errorf("could not find executable: \"%s\": %w", commandParts[0], err)
	}

	// If the executable is a symlink, resolve it.
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return nil, fmt.Errorf("could not resolve executable: \"%s\": %w", commandParts[0], err)
	}

	executableDetails, err := os.Stat(executable)
	if err != nil {
		return nil, fmt.Errorf("could not stat executable:  \"%s\": %w", commandParts[0], err)
	}

	var cacheKey string

	// Try to get the value from the cache.
	if getter.Cache.Enabled {
		cacheKey = fmt.Sprintf(
			"%s %s -- %d/%d",
			executable,
			strings.Join(commandParts[1:], " "),
			executableDetails.ModTime().Unix(),
			executableDetails.Size(),
		)

		if value := valueCache.Get(cacheKey); value != nil {
			return value, nil
		}
	}

	// If that fails, run the command.
	cmd := exec.Command(executable, commandParts[1:]...)
	cmd.Dir = projectFolder.Path()
	result, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running command: \"%s\": %w", executable, err)
	}

	// Store the value in the cache for future generations.
	if cacheKey != "" {
		valueCache.Set(cacheKey, result)
	}

	return result, nil
}

func (getter CustomGetter) getAncestorFileValue(
	projectFolder fileutils.Directory,
	filePath string,
) ([]byte, error) {
	resolvedFilePath := projectFolder.FindFileInAncestors(filePath)
	if resolvedFilePath == "" {
		return nil, fmt.Errorf("could not find file: \"%s\"", filePath)
	}

	return os.ReadFile(resolvedFilePath)
}

func (getter CustomGetter) getValueAs(
	as AsType,
	value []byte,
) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	switch as {
	case AsJSON:
		err := json.Unmarshal(value, &result)
		if err != nil {
			return nil, fmt.Errorf("invalid json: \"%s\": %w", value, err)
		}
	case AsYAML:
		err := yaml.Unmarshal(value, &result)
		if err != nil {
			return nil, fmt.Errorf("invalid yaml: \"%s\": %w", value, err)
		}
	case AsTOML:
		err := toml.Unmarshal(value, &result)
		if err != nil {
			return nil, fmt.Errorf("invalid toml: \"%s\": %w", value, err)
		}
	case AsText:
		result = map[string]interface{}{
			"Text": strings.TrimSpace(string(value)),
		}

	default:
		return nil, fmt.Errorf("invalid value type: \"%s\"", as)
	}

	return result, nil
}

func (getter CustomGetter) applyTemplate(as AsType, bytesValue []byte) (interface{}, error) {
	var result interface{}

	// Parse the value into a map.
	result, err := getter.getValueAs(as, bytesValue)
	if err != nil {
		return "", err
	}

	if getter.ValueTemplate != "" {
		// Run the value through the ValueTemplate.
		tmpl := template.New(getter.Type.String()).Funcs(sprigTemplateFunctions)
		tmpl, err = tmpl.Parse(getter.ValueTemplate)
		if err != nil {
			return "", fmt.Errorf("invalid template: \"%s\": %w", getter.ValueTemplate, err)
		}

		result, err = modtemplate.TemplateToString(tmpl, result)
		if err != nil {
			return "", fmt.Errorf("error executing template: \"%s\": %w", getter.ValueTemplate, err)
		}
	}

	return result, nil
}
