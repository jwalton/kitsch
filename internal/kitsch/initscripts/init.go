package initscripts

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/init*
var initTemplates embed.FS

func getKitschCommand() string {
	return "kitsch"
}

// ShortInitScript returns the kitsch initialization script for the given shell type.
func ShortInitScript(shell string, configFile string) (string, error) {
	return getInitScript("templates/init-short.", shell, configFile)
}

// InitScript returns the full kitsch initialization script for the given shell type.
func InitScript(shell string, configFile string) (string, error) {
	return getInitScript("templates/init.", shell, configFile)
}

func getInitScript(filename string, shell string, configFile string) (string, error) {
	data := map[string]string{
		"kitschCommand": getKitschCommand(),
		"configFile":    configFile,
	}

	initTemplate, err := initTemplates.ReadFile(filename + shell)
	if err != nil {
		return "", fmt.Errorf("Invalid shell %s", shell)
	}

	return execTemplate(string(initTemplate), data)
}

func execTemplate(templateSrc string, data interface{}) (string, error) {
	t := template.Must(template.New("template").Parse(templateSrc))

	var b bytes.Buffer
	err := t.Execute(&b, data)

	return b.String(), err
}
