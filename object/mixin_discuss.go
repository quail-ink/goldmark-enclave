package object

import (
	"bytes"
	"text/template"
)

const discussWidgetTpl = `
	<video
  controls
  src="{{.URL}}"
  width="640"
  height="480"
  >
		Sorry, your browser does not support embedded videos.
	</video>
`

func GetDiscussWidgetHtml(url string) (string, error) {
	t, err := template.New("discuss-widget").Parse(discussWidgetTpl)
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
