package object

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"github.com/quail-ink/goldmark-enclave/core"
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

func GetQuailWidgetHtml(enc *core.Enclave) (string, error) {
	if enc.Theme == "dark" {
		enc.Theme = "dark"
	} else {
		enc.Theme = "light"
	}
	var err error

	ret := ""
	buf := bytes.Buffer{}
	if enc.IframeDisabled {
		ret, err = GetNoIframeTplHtml(enc, fmt.Sprintf("%s://%s%s", enc.URL.Scheme, enc.URL.Host, enc.URL.Path))
		if err != nil {
			return "", err
		}

	} else {
		t, err := template.New("quail-widget").Parse(quailWidgetTpl)
		if err != nil {
			return "", err
		}

		layout := ""
		if l, ok := enc.Params["layout"]; ok {
			layout = l
		}

		height := "auto"
		if strings.Contains(enc.URL.Path, "/p/") {
			height = "128px"
		} else if layout == "subscribe_form" {
			height = "390px"
		} else if layout == "subscribe_form_mini" {
			height = "142px"
		}

		if err = t.Execute(&buf, map[string]string{
			"URL":    fmt.Sprintf("%s://%s%s/widget?theme=%s&layout=%s&logged=ignore", enc.URL.Scheme, enc.URL.Host, enc.URL.Path, enc.Theme, layout),
			"Theme":  enc.Theme,
			"Height": height,
		}); err != nil {
			return "", err
		}

		ret = buf.String()
	}

	return ret, nil
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
