package modtemplate

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/env"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/powerline"
	"github.com/jwalton/kitsch-prompt/internal/kitsch/styling"
)

const recursionMaxNums = 1000

var sprigTemplateFunctions template.FuncMap

func init() {
	sprigTemplateFunctions = sprig.TxtFuncMap()
	delete(sprigTemplateFunctions, "env")
}

// CompileTemplate compiles a module template and adds default template functions.
func CompileTemplate(styles *styling.Registry, environment env.Env, name string, templateString string) (*template.Template, error) {
	tmpl := template.New(name)

	funcMap := template.FuncMap{}

	// Borrowed from Helm.
	// TODO: Make `data` optional.
	includedNames := make(map[string]int) // Recursion guard.
	funcMap["include"] = func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > recursionMaxNums {
				return "", fmt.Errorf("rendering template has a nested reference name: %s", name)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := tmpl.ExecuteTemplate(&buf, name, data)
		includedNames[name]--
		return buf.String(), err
	}

	funcMap["env"] = func(name string) string {
		return environment.Getenv(name)
	}

	tmpl, err := tmpl.
		Funcs(funcMap).
		Funcs(sprigTemplateFunctions).
		Funcs(styling.TxtFuncMap(styles)).
		Funcs(powerline.TxtFuncMap(styles)).
		Parse(templateString)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// TemplateToString renders a template to a string.
func TemplateToString(template *template.Template, data interface{}) (string, error) {
	var b bytes.Buffer
	err := template.Execute(&b, data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
