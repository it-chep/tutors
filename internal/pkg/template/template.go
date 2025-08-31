package template

import (
	"bytes"
	"text/template"
	"time"
)

func Format(t time.Time) string {
	return t.Format(time.DateOnly)
}

var funcMap = map[string]any{
	"date_format": Format,
}

// Execute исполняет шаблон и в случае ошибки отдает пустую строку
func Execute(tmpl string, value any) string {
	t, err := template.New("").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return ""
	}

	var buff bytes.Buffer
	if err = t.Execute(&buff, value); err != nil {
		return ""
	}
	return buff.String()
}
