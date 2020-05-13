package render

import (
	"bytes"
	"html/template"
	texttpl "text/template"
)

func HtmlRender(name, tmpl string, data interface{}) (string, error) {
	var buf bytes.Buffer
	tpl := template.New(name)
	tpl.Funcs(template.FuncMap(defaultFuncMaps))
	t, err := tpl.Parse(tmpl)
	if err != nil {
		return "", err
	}

	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TextRender(name, tmpl string, data interface{}) (string, error) {
	var buf bytes.Buffer
	tpl := texttpl.New(name)
	tpl.Funcs(defaultFuncMaps)
	t, err := tpl.Parse(tmpl)
	if err != nil {
		return "", err
	}

	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
