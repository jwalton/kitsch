package getters

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig/v3"
	"github.com/jwalton/kitsch-prompt/internal/fileutils"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/modtemplate"
	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v3"
)

var sprigTemplateFunctions = sprig.TxtFuncMap()

// CustomGetter is a getter that can be configured from a YAML file.
type CustomGetter struct {
	// Type is the type of getter.  One of "custom", "file", "ancestorFile", or "env".
	Type string `yaml:"type"`
	// From is the source to get data from.  The meaning of "From" is based on
	// the provided "Type".
	From string `yaml:"from"`
	// As will determine how to interpret the result of the getter.  One of "text", "json", "toml", or "yaml".
	As string `yaml:"as"`
	// ValueTemplate is a golang template used to parse values out of the result of
	// the getter.
	ValueTemplate string `yaml:"valueTemplate"`
	// Regex is a regular expression used to parse values out of the result of
	// the getter.  If specified, then "As" and "Template" will be ignored.
	Regex string `yaml:"regex"`
}

// GetValue gets the value for this getter.  The return value will be either a string,
// of if the value is a JSON, YAML, or TOML object, and the `ValueTemplate` is not set,
// the parsed contents of the object.
func (getter CustomGetter) GetValue(
	folder fileutils.Directory,
) (interface{}, error) {
	// Get the raw value for the getter.
	var bytesValue []byte
	var err error
	var result interface{}

	switch getter.Type {
	case "custom":
		bytesValue, err = getter.getCustomValue(folder, getter.From)
	case "file":
		bytesValue, err = fs.ReadFile(folder.FileSystem(), getter.From)
	case "ancestorFile":
		bytesValue, err = getter.getAncestorFileValue(folder, getter.From)
	case "env":
		strValue := os.Getenv(getter.From)
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
			result, err = getter.applyTemplate("text", []byte(strValue))
			if err != nil {
				return "", err
			}
		} else {
			result = strValue
		}

	} else if (getter.As != "" || getter.ValueTemplate != "") && !(getter.As == "text" && getter.ValueTemplate == "") {
		as := getter.As
		if as == "" {
			as = "text"
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

	cmd := exec.Command(executable, commandParts[1:]...)
	cmd.Dir = projectFolder.Path()
	return cmd.CombinedOutput()
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
	as string,
	value []byte,
) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	switch as {
	case "json":
		err := json.Unmarshal(value, &result)
		if err != nil {
			return nil, fmt.Errorf("invalid json: \"%s\": %w", value, err)
		}
	case "yaml":
		err := yaml.Unmarshal(value, &result)
		if err != nil {
			return nil, fmt.Errorf("invalid yaml: \"%s\": %w", value, err)
		}
	case "toml":
		err := toml.Unmarshal(value, &result)
		if err != nil {
			return nil, fmt.Errorf("invalid toml: \"%s\": %w", value, err)
		}
	case "text":
		result = map[string]interface{}{
			"Text": strings.TrimSpace(string(value)),
		}

	default:
		return nil, fmt.Errorf("invalid value type: \"%s\"", as)
	}

	return result, nil
}

func (getter CustomGetter) applyTemplate(as string, bytesValue []byte) (interface{}, error) {
	var result interface{}

	// Parse the value into a map.
	result, err := getter.getValueAs(as, bytesValue)
	if err != nil {
		return "", err
	}

	if getter.ValueTemplate != "" {
		// Run the value through the ValueTemplate.
		tmpl := template.New(getter.Type).Funcs(sprigTemplateFunctions)
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
