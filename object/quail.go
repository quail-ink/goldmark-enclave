package object

import (
	"bytes"
	"fmt"
	"net/url"
	"text/template"
)

const quailTpl = `
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

func GetQuailWidgetHtml(url *url.URL, theme string) (string, error) {
	if theme == "dark" {
		theme = "dark"
	} else {
		theme = "light"
	}

	t, err := template.New("quail").Parse(quailTpl)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	if err = t.Execute(&buf, map[string]string{
		"URL":   fmt.Sprintf("%s://%s%s/widget#theme=%s", url.Scheme, url.Host, url.Path, theme),
		"Theme": theme,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
