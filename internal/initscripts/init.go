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

// InitScript returns the kitsch-prompt initalization script for the given
// shell type.
func InitScript(shell string) (string, error) {
	kitschcommand, err := os.Executable()
	if err != nil {
		kitschcommand = "kitsch-prompt"
	}

	data := map[string]string{
		"kitschcommand": kitschcommand,
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
