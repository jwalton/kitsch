package initscripts

import (
	"bytes"
	// embed required for templates below.
	_ "embed"
	"os"
	"text/template"
)

//go:embed templates/init.zsh
var zshTemplate string

// InitScript returns the kitsch-prompt initialization script for the given
// shell type.
func InitScript(shell string, configFile string) (string, error) {
	kitschCommand, err := os.Executable()
	if err != nil {
		kitschCommand = "kitsch-prompt"
	}

	data := map[string]string{
		"kitschCommand": kitschCommand,
		"configFile":    configFile,
	}

	switch shell {
	case "zsh":
		return execTemplate(zshTemplate, data)
	default:
		panic("Invalid shell type: " + shell)
	}

}

func execTemplate(templateSrc string, data interface{}) (string, error) {
	t := template.Must(template.New("template").Parse(templateSrc))

	var b bytes.Buffer
	err := t.Execute(&b, data)

	return b.String(), err
}
