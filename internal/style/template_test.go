package style

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/jwalton/gchalk"
	"github.com/stretchr/testify/assert"
)

// Compilte a module template and add default template functions.
func testCompileTemplate(name string, templateString string) *template.Template {
	tmpl := template.Must(template.New(name).Funcs(TxtFuncMap()).Parse(templateString))
	return tmpl
}

func testTemplateToString(template *template.Template, data interface{}) string {
	var b bytes.Buffer
	err := template.Execute(&b, data)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

func TestStyleFunc(t *testing.T) {
	gchalk.SetLevel(gchalk.LevelAnsi16m)

	tmpl := testCompileTemplate("test", `{{ . | style "red" }}`)
	assert.Equal(t, "\u001B[31mfoo\u001B[39m", testTemplateToString(tmpl, "foo"))

	tmpl = testCompileTemplate("test", `{{ . | style "red" "bgBlue"}}`)
	assert.Equal(t, "\u001B[44m\u001B[31mfoo\u001B[39m\u001B[49m", testTemplateToString(tmpl, "foo"))

	tmpl = testCompileTemplate("test", `{{ . | style "red bgBlue"}}`)
	assert.Equal(t, "\u001B[44m\u001B[31mfoo\u001B[39m\u001B[49m", testTemplateToString(tmpl, "foo"))

	// Should work with no styles.
	tmpl = testCompileTemplate("test", `{{ style . }}`)
	assert.Equal(t, "foo", testTemplateToString(tmpl, "foo"))

	// Should not crash if there's no arguments at all.
	tmpl = testCompileTemplate("test", `{{ style }}`)
	assert.Equal(t, "", testTemplateToString(tmpl, "foo"))
}

func TestFgColorFunc(t *testing.T) {
	gchalk.SetLevel(gchalk.LevelAnsi16m)

	tmpl := testCompileTemplate("test", `{{ . | fgColor "red" }}`)
	assert.Equal(t, "\u001B[31mfoo\u001B[39m", testTemplateToString(tmpl, "foo"))

	tmpl2 := testCompileTemplate("test", `{{ . | fgColor "bgRed"}}`)
	assert.Equal(t, "\u001B[31mfoo\u001B[39m", testTemplateToString(tmpl2, "foo"))

	tmpl3 := testCompileTemplate("test", `{{ . | fgColor "bg:red"}}`)
	assert.Equal(t, "\u001B[31mfoo\u001B[39m", testTemplateToString(tmpl3, "foo"))
}

func TestBgColorFunc(t *testing.T) {
	gchalk.SetLevel(gchalk.LevelAnsi16m)

	tmpl := testCompileTemplate("test", `{{ . | bgColor "red" }}`)
	assert.Equal(t, "\u001B[41mfoo\u001B[49m", testTemplateToString(tmpl, "foo"))

	tmpl2 := testCompileTemplate("test", `{{ . | bgColor "bgRed"}}`)
	assert.Equal(t, "\u001B[41mfoo\u001B[49m", testTemplateToString(tmpl2, "foo"))

	tmpl3 := testCompileTemplate("test", `{{ . | bgColor "bg:red"}}`)
	assert.Equal(t, "\u001B[41mfoo\u001B[49m", testTemplateToString(tmpl3, "foo"))
}
