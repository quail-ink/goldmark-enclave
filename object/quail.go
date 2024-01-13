package object

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"text/template"
)

const quailWidgetTpl = `
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

const quailImageTpl = `
<div class="quail-image-wrapper" style="width: {{.Width}}; height: {{.Height}}">
	<img src="{{.URL}}" alt="{{.Title}}" style="width: {{.Width}}; height: {{.Height}}" class="quail-image" />
	<div class="quail-image-caption">{{.Title}}</div>
</div>
`

func GetQuailWidgetHtml(url *url.URL, theme string) (string, error) {
	if theme == "dark" {
		theme = "dark"
	} else {
		theme = "light"
	}

	t, err := template.New("quail-widget").Parse(quailWidgetTpl)
	if err != nil {
		return "", err
	}

	height := 96
	if strings.Contains(url.Path, "/p/") {
		height = 128
	}

	buf := bytes.Buffer{}
	if err = t.Execute(&buf, map[string]string{
		"URL":    fmt.Sprintf("%s://%s%s/widget?theme=%s", url.Scheme, url.Host, url.Path, theme),
		"Theme":  theme,
		"Height": strconv.Itoa(height),
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetQuailImageHtml(url *url.URL, params map[string]string) (string, error) {
	buf := bytes.Buffer{}

	t, err := template.New("quail-image").Parse(quailImageTpl)
	if err != nil {
		return "", err
	}

	w := "auto"
	if width, ok := params["width"]; ok {
		w = width
	}

	h := "auto"
	if height, ok := params["height"]; ok {
		h = height
	}
	title := ""
	if t, ok := params["title"]; ok {
		title = t
	}

	if err = t.Execute(&buf, map[string]string{
		"URL":    url.String(),
		"Title":  title,
		"Width":  w,
		"Height": h,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
