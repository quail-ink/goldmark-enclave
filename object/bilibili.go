package object

import (
	"bytes"
	"text/template"

	"github.com/quailyquaily/goldmark-enclave/core"
)

const bilibiliTpl = `
<iframe class="enclave-object bilibili-enclave-object"
	width="100%" height="400"
	src="https://player.bilibili.com/player.html?bvid={{.ObjectID}}&page=1&autoplay=0"
	scrolling="no" border="0" framespacing="0" allowfullscreen="true" frameborder="no">
</iframe>
`

func GetBilibiliEmbedHtml(enc *core.Enclave) (string, error) {
	var err error
	ret := ""
	if enc.IframeDisabled {
		ret, err = GetNoIframeTplHtml(enc, string(enc.Image.Destination))
		if err != nil {
			return "", err
		}

	} else {
		t, err := template.New("bilibili").Parse(bilibiliTpl)
		if err != nil {
			return "", err
		}

		buf := bytes.Buffer{}
		if err = t.Execute(&buf, map[string]string{
			"ObjectID": enc.ObjectID,
		}); err != nil {
			return "", err
		}
		ret = buf.String()
	}

	return ret, nil
}
