package enclave

import (
	"fmt"

	"github.com/quail-ink/goldmark-enclave/object"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type HTMLRenderer struct {
	cfg *Config
}

func NewHTMLRenderer(cfg *Config) renderer.NodeRenderer {
	r := &HTMLRenderer{cfg: cfg}
	return r
}

func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// image with alt like [alt](url "title") will generate a node seq like
	// layout:
	// - imgLeftNode: kind = paragraph, content = alt
	// - imgNode: kind = image
	// - imgRightNode: kind = text, content = alt
	// I don't know how to handle them yet.
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
	case EnclaveProviderQuailWidget:
		html, err := object.GetQuailWidgetHtml(enc.URL, enc.Theme, enc.Params)
		if err != nil || html == "" {
			html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object quail-enclave-object error">Failed to load quail widget from %s</div></div>`, enc.ObjectID)
		} else {
			html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object quail-enclave-object normal-object no-border">%s</div></div>`, html)
		}
		w.Write([]byte(html))
	case EnclaveProviderQuailImage:
		alt := string(node.Text(source))
		if alt == "" && len(enc.Title) != 0 {
			alt = fmt.Sprintf("An image to describe %s", enc.Title)
		}
		if alt == "" {
			alt = "An image to describe post"
		}
		enc.Params["alt"] = alt
		html, err := object.GetQuailImageHtml(enc.URL, enc.Params)
		if err != nil || html == "" {
			html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object quail-enclave-object error">Failed to load quail image from %s</div></div>`, enc.ObjectID)
		}
		w.Write([]byte(html))
	case EnclaveRegularImage:
		alt := string(node.Text(source))
		if alt == "" && len(enc.Title) != 0 {
			alt = fmt.Sprintf("An image to describe %s", enc.Title)
		}
		if alt == "" {
			alt = "An image to describe post"
		}
		html := fmt.Sprintf(`<img src="%s" alt="%s" />`, enc.URL.String(), alt)
		fmt.Printf("regular image html: %+v\n", html)
		w.Write([]byte(html))
	}

	return ast.WalkContinue, nil
}
