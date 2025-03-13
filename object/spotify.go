package object

import (
	"bytes"
	"text/template"

	"github.com/quailyquaily/goldmark-enclave/core"
)

const SpotifyTpl = `
<iframe style="border-radius:12px"
	src="https://open.spotify.com/embed/track/{{.TrackID}}?utm_source=generator"
	width="100%" height="352" frameBorder="0" allowfullscreen=""
	allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture" loading="lazy">
</iframe>
`

func GetSpotifyWidgetHtml(enc *core.Enclave) (string, error) {
	var err error
	ret := ""
	if enc.IframeDisabled {
		ret, err = GetNoIframeTplHtml(enc, string(enc.Image.Destination))
		if err != nil {
			return "", err
		}

	} else {
		t, err := template.New("spotify").Parse(SpotifyTpl)
		if err != nil {
			return "", err
		}

		buf := bytes.Buffer{}
		if err = t.Execute(&buf, map[string]string{
			"TrackID": enc.ObjectID,
		}); err != nil {
			return "", err
		}
		ret = buf.String()
	}

	return ret, nil
}
