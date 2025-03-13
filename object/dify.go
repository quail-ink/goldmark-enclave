package object

import (
	"bytes"
	"text/template"

	"github.com/quailyquaily/goldmark-enclave/core"
)

const difyWidgetTpl = `
	<iframe
		src="{{.URL}}"
		style="width: 100%; height: 100%; min-height: 700px"
		frameborder="0"
		allow="microphone">
	</iframe>
`

func GetDifyWidgetHtml(enc *core.Enclave) (string, error) {
	var err error
	ret := ""
	if enc.IframeDisabled {
		ret, err = GetNoIframeTplHtml(enc, enc.ObjectID)
		if err != nil {
			return "", err
		}

	} else {
		t, err := template.New("dify-widget").Parse(difyWidgetTpl)
		if err != nil {
			return "", err
		}

		buf := bytes.Buffer{}
		if err = t.Execute(&buf, map[string]string{
			"URL": enc.ObjectID,
		}); err != nil {
			return "", err
		}
		ret = buf.String()
	}

	return ret, nil
}
