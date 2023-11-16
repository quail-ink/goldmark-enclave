package enclave

import (
	"fmt"

	"github.com/quail-ink/goldmark-enclave/object"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type HTMLRenderer struct{}

func NewHTMLRenderer() renderer.NodeRenderer {
	r := &HTMLRenderer{}
	return r
}

func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindEnclave, r.renderEnclave)
}

func (r *HTMLRenderer) renderEnclave(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkContinue, nil
	}

	enc := node.(*Enclave)
	switch enc.Provider {
	case EnclaveProviderYouTube:
		w.Write([]byte(`<div class="enclave-object-wrapper"><iframe class="enclave-object youtube-enclave-object" width="100%" height="400" src="https://www.youtube.com/embed/` + enc.ObjectID + `" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe></div>`))
	case EnclaveProviderBilibili:
		w.Write([]byte(`<div class="enclave-object-wrapper"><iframe class="enclave-object bilibili-enclave-object" width="100%" height="400" src="//player.bilibili.com/player.html?bvid=` + enc.ObjectID + `&page=1" scrolling="no" border="0" framespacing="0" allowfullscreen="true" frameborder="no"></iframe></div>`))
	case EnclaveProviderTwitter:
		html, err := object.GetTweetOembedHtml(enc.ObjectID, enc.Theme)
		if err != nil || html == "" {
			html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object twitter-enclave-object normal-object error">Failed to load tweet from %s</div></div>`, enc.ObjectID)
		} else {
			html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object twitter-enclave-object normal-object no-border">%s</div></div>`, html)
		}
		w.Write([]byte(html))
	case EnclaveProviderTradingView:
		html, err := object.GetTradingViewWidgetHtml(enc.ObjectID, enc.Theme)
		if err != nil || html == "" {
			html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object tradingview-enclave-object error">Failed to load tradingview chart from %s</div></div>`, enc.ObjectID)
		} else {
			html = fmt.Sprintf(`<div class="enclave-object-wrapper auto-resize"><div class="enclave-object tradingview-enclave-object no-border">%s</div></div>`, html)
		}
		w.Write([]byte(html))
	case EnclaveProviderQuail:
		html, err := object.GetQuailWidgetHtml(enc.ObjectID, enc.Theme)
		if err != nil || html == "" {
			html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object quail-enclave-object error">Failed to load quail widget from %s</div></div>`, enc.ObjectID)
		} else {
			html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object quail-enclave-object normal-object">%s</div></div>`, html)
		}
		w.Write([]byte(html))
	}

	return ast.WalkContinue, nil
}
