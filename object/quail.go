package object

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"text/template"
)

const quailWidgetTpl = `
<iframe
	src="{{.URL}}"
	data-theme="{{.Theme}}"
	width="100%"
	height="{{.Height}}"
	title="Quail Widget"
	frameborder="0"
	allow="web-share"
	allowfullscreen
></iframe>
`

const quailImageTpl = `
<span class="quail-image-wrapper" style="width: {{.Width}}; height: {{.Height}}; margin: 1rem 0; display: block">
	<img src="{{.URL}}" alt="{{.Alt}}" style="width: {{.Width}}; height: {{.Height}}" class="quail-image" />
	<span class="quail-image-caption" style="display: block">{{.Title}}</span>
</span>
`

func GetQuailWidgetHtml(url *url.URL, theme string, params map[string]string) (string, error) {
	if theme == "dark" {
		theme = "dark"
	} else {
		theme = "light"
	}

	t, err := template.New("quail-widget").Parse(quailWidgetTpl)
	if err != nil {
		return "", err
	}

	layout := ""
	if l, ok := params["layout"]; ok {
		layout = l
	}

	height := "auto"
	if strings.Contains(url.Path, "/p/") {
		height = "128px"
	} else if layout == "subscribe_form" {
		height = "390px"
	} else if layout == "subscribe_form_mini" {
		height = "142px"
	}

	buf := bytes.Buffer{}
	if err = t.Execute(&buf, map[string]string{
		"URL":    fmt.Sprintf("%s://%s%s/widget?theme=%s&layout=%s&logged=ignore", url.Scheme, url.Host, url.Path, theme, layout),
		"Theme":  theme,
		"Height": height,
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

	alt := ""
	if t, ok := params["alt"]; ok {
		alt = t
	}

	if err = t.Execute(&buf, map[string]string{
		"URL":    url.String(),
		"Title":  title,
		"Width":  w,
		"Height": h,
		"Alt":    alt,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
