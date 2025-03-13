package object

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/quailyquaily/goldmark-enclave/core"
)

const youtubeTpl = `
<iframe
	class="enclave-object youtube-enclave-object"
	width="100%" height="400"
	src="https://www.youtube.com/embed/{{.ObjectID}}"
	title="YouTube video player"
	frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen>
</iframe>
`

func GetYoutubeEmbedHtml(enc *core.Enclave) (string, error) {
	var err error
	ret := ""
	if enc.IframeDisabled {
		coverImageUrl := fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", enc.ObjectID)
		ret, err = GetNoIframeTplHtmlWithAttrs(enc, string(enc.Image.Destination), coverImageUrl, "YouTube video player")
		if err != nil {
			return "", err
		}

	} else {
		t, err := template.New("youtube").Parse(youtubeTpl)
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
