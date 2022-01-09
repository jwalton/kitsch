package initscripts

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed templates/*init*
var initTemplates embed.FS

func getKitschCommand() string {
	kitschCommand, err := os.Executable()
	if err != nil {
		kitschCommand = "kitsch"
	}

	return kitschCommand
}

// ShortInitScript returns the kitsch initialization script for the given shell type.
func ShortInitScript(shell string, configFile string) (string, error) {
	return getInitScript("init-short", shell, configFile)
}

// InitScript returns the full kitsch initialization script for the given shell type.
func InitScript(shell string, configFile string) (string, error) {
	return getInitScript("init", shell, configFile)
}

func getInitScript(filename string, shell string, configFile string) (string, error) {
	kitschCommand := getKitschCommand()

	shellExt := shell
	if shell == "powershell" {
		kitschCommand = `"` + kitschCommand + `"`
		shellExt = "ps1"
	}

	data := map[string]string{
		"kitschCommand": kitschCommand,
		"configFile":    configFile,
	}

	initTemplate, err := initTemplates.ReadFile("templates/" + shell + "-" + filename + "." + shellExt)
	if err != nil {
		return "", fmt.Errorf("invalid shell %s", shell)
	}

	return execTemplate(string(initTemplate), data)
}

func execTemplate(templateSrc string, data interface{}) (string, error) {
	t := template.Must(template.New("template").Parse(templateSrc))

	var b bytes.Buffer
	err := t.Execute(&b, data)

	return b.String(), err
}
