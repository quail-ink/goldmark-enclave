package object

import (
	"bytes"
	"text/template"

	"github.com/quail-ink/goldmark-enclave/core"
)

const audioTpl = `<audio controls src="{{.URL}}"></audio>`

func GetAudioHtml(enc *core.Enclave) (string, error) {
	var err error
	ret := ""

	t, err := template.New("audio").Parse(audioTpl)
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

	return ret, nil
}
