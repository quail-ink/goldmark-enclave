package object

import (
	"bytes"
	"text/template"

	"github.com/quailyquaily/goldmark-enclave/core"
)

const audioTpl = `<div class="enclave-audio-wrapper flex place-center">
<audio controls src="{{.URL}}"></audio>
</div>`

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
