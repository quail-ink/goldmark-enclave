package object

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"text/template"
)

const quailTpl = `
<iframe
	src="{{.URL}}"
	data-theme="{{.Theme}}"
	width="100%"
	height="{{.Height}}px"
	title="Quail Widget"
	frameborder="0"
	allow="web-share"
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

	height := 96
	if strings.Contains(url.Path, "/p/") {
		height = 128
	}

	buf := bytes.Buffer{}
	if err = t.Execute(&buf, map[string]string{
		"URL":    fmt.Sprintf("%s://%s%s/widget#theme=%s", url.Scheme, url.Host, url.Path, theme),
		"Theme":  theme,
		"Height": strconv.Itoa(height),
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
