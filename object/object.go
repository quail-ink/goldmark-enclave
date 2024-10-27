package object

import (
	"bytes"
	"text/template"

	"github.com/quail-ink/goldmark-enclave/core"
)

const NoIframeTpl = `
<a href="{{.URL}}" target="_blank" rel="noopener noreferrer">
	<img src="{{.Src}}" />
</a>
`

func GetNoIframeTplHtml(enc *core.Enclave, url string) (string, error) {
	buf := bytes.Buffer{}
	if enc.IframeDisabled {
		t, err := template.New("no-iframe").Parse(NoIframeTpl)
		if err != nil {
			return "", err
		}

		if err = t.Execute(&buf, map[string]string{
			"URL":   url,
			"Src": core.IframeDisabledPlaceholderURL,
		}); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}
