package object

import (
	"bytes"
	"fmt"
	"text/template"
)

const quailListTpl = `
<iframe
	src="{{.URL}}"
	data-theme="{{.Theme}}"
	width="100%"
	height="96px"
	title="Quail Widget"
	frameborder="0"
	allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; web-share"
	allowfullscreen
></iframe>
`

func GetQuailWidgetHtml(url, theme string) (string, error) {
	if theme == "dark" {
		theme = "dark"
	} else {
		theme = "light"
	}

	t, err := template.New("quail-list").Parse(quailListTpl)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	if err = t.Execute(&buf, map[string]string{
		"URL":   fmt.Sprintf("%s/widget", url),
		"Theme": theme,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
