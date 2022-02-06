package modtemplate

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/jwalton/go-ansiparser"
	"github.com/jwalton/kitsch/internal/kitsch/env"
	"github.com/jwalton/kitsch/internal/kitsch/powerline"
	"github.com/jwalton/kitsch/internal/kitsch/styling"
)

const csi = "\x1b["

func addTemplateFunctions(
	styles *styling.Registry,
	environment env.Env,
	screenWidth int,
	tmpl *template.Template,
) *template.Template {
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

	funcMap["rightJustify"] = func(str string) string {
		strs := strings.Split(str, "\n")
		for i, s := range strs {
			printLength := 0
			ansiTokenizer := ansiparser.NewStringTokenizer(s)
			for ansiTokenizer.Next() {
				printLength += ansiTokenizer.Token().PrintLength()
			}
			strs[i] = fmt.Sprintf("%ss%s%dG%s%su", csi, csi, screenWidth-printLength, s, csi)
		}
		return strings.Join(strs, "\n")
	}

	return tmpl.Funcs(funcMap).
		Funcs(sprigTemplateFunctions).
		Funcs(styling.TxtFuncMap(styles)).
		Funcs(powerline.TxtFuncMap(styles))
}
