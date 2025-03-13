package object

import (
	"bytes"
	"text/template"

	"github.com/quailyquaily/goldmark-enclave/core"
)

const NoIframeTpl = `
<a href="{{.URL}}" target="_blank" rel="noopener noreferrer" title="{{.Title}}">
	<img src="{{.Src}}" alt="{{.Title}}"/>
</a>
`

func GetNoIframeTplHtml(enc *core.Enclave, url string) (string, error) {
	return GetNoIframeTplHtmlWithAttrs(enc, url, core.IframeDisabledPlaceholderURL, "Click to view")
}

func GetNoIframeTplHtmlWithAttrs(enc *core.Enclave, url, src, title string) (string, error) {
	buf := bytes.Buffer{}
	t, err := template.New("no-iframe").Parse(NoIframeTpl)
	if err != nil {
		return "", err
	}

	if err = t.Execute(&buf, map[string]string{
		"URL":   src,
		"Src":   src,
		"Title": title,
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
