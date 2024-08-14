package object

import (
	"bytes"
	"text/template"
)

const difyWidgetTpl = `
	<iframe
		src="{{.URL}}"
		style="width: 100%; height: 100%; min-height: 700px"
		frameborder="0"
		allow="microphone">
	</iframe>
`

func GetDifyWidgetHtml(url string) (string, error) {
	t, err := template.New("dify-widget").Parse(difyWidgetTpl)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	if err = t.Execute(&buf, map[string]string{
		"URL": url,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
