package enclave

import (
	"fmt"

	"github.com/quail-ink/goldmark-enclave/core"
	"github.com/quail-ink/goldmark-enclave/object"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type HTMLRenderer struct {
	cfg *core.Config
}

func NewHTMLRenderer(cfg *core.Config) renderer.NodeRenderer {
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
	reg.Register(core.KindEnclave, r.renderEnclave)
}

func (r *HTMLRenderer) renderEnclave(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkContinue, nil
	}

	enc := node.(*core.Enclave)
	switch enc.Provider {
	case core.EnclaveProviderYouTube:
		{
			// if enc.IframeDisabled {
			// 	w.Write([]byte(`
			// 	<div class="enclave-object-wrapper">
			// 		<a href="https://www.youtube.com/watch?v=` + enc.ObjectID + `" target="_blank" rel="noopener noreferrer">
			// 			<img src="https://img.youtube.com/vi/` + enc.ObjectID + `/maxresdefault.jpg" title="YouTube video player" />
			// 		</a>
			// 	</div>`))
			// } else {
			// 	w.Write([]byte(`
			// 	<div class="enclave-object-wrapper">
			// 		<iframe class="enclave-object youtube-enclave-object" width="100%" height="400" src="https://www.youtube.com/embed/` + enc.ObjectID + `" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>
			// 	</div>`))
			// }
			html, err := object.GetYoutubeEmbedHtml(enc)
			if err != nil || html == "" {
				html = r.wrapEnclaveErrorHtml("youtube", enc.ObjectID)
			} else {
				html = r.wrapEnclaveHtml("youtube", html, false, false)
			}
			w.Write([]byte(html))
		}

	case core.EnclaveProviderBilibili:
		{
			html, err := object.GetBilibiliEmbedHtml(enc)
			if err != nil || html == "" {
				html = r.wrapEnclaveErrorHtml("bilibili", enc.ObjectID)
			} else {
				html = r.wrapEnclaveHtml("bilibili", html, false, false)
			}
			w.Write([]byte(html))
		}

	case core.EnclaveProviderTwitter:
		html, err := object.GetTweetOembedHtml(enc.ObjectID, enc.Theme)
		if err != nil || html == "" {
			// html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object twitter-enclave-object normal-object error">Failed to load tweet from %s</div></div>`, enc.ObjectID)
			html = r.wrapEnclaveErrorHtml("twitter", enc.ObjectID)
		} else {
			// html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object twitter-enclave-object normal-object no-border">%s</div></div>`, html)
			html = r.wrapEnclaveHtml("twitter", html, true, false)
		}
		w.Write([]byte(html))

	case core.EnclaveProviderTradingView:
		html, err := object.GetTradingViewWidgetHtml(enc)
		if err != nil || html == "" {
			// html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object tradingview-enclave-object error">Failed to load tradingview chart from %s</div></div>`, enc.ObjectID)
			html = r.wrapEnclaveErrorHtml("tradingview", enc.ObjectID)
		} else {
			// html = fmt.Sprintf(`<div class="enclave-object-wrapper auto-resize"><div class="enclave-object tradingview-enclave-object no-border">%s</div></div>`, html)
			html = r.wrapEnclaveHtml("tradingview", html, false, false)
		}
		w.Write([]byte(html))

	case core.EnclaveProviderDifyWidget:
		html, err := object.GetDifyWidgetHtml(enc)
		if err != nil || html == "" {
			// html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object dify-enclave-object error">Failed to load dify widget from %s</div></div>`, enc.ObjectID)
			html = r.wrapEnclaveErrorHtml("dify", enc.ObjectID)
		} else {
			// html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object dify-enclave-object normal-object no-border">%s</div></div>`, html)
			html = r.wrapEnclaveHtml("dify", html, true, false)
		}
		w.Write([]byte(html))

	case core.EnclaveProviderQuailWidget:
		html, err := object.GetQuailWidgetHtml(enc)
		if err != nil || html == "" {
			// html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object quail-enclave-object error">Failed to load quail widget from %s</div></div>`, enc.ObjectID)
			html = r.wrapEnclaveErrorHtml("quail", enc.ObjectID)
		} else {
			// html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object quail-enclave-object normal-object no-border">%s</div></div>`, html)
			html = r.wrapEnclaveHtml("quail", html, true, false)
		}
		w.Write([]byte(html))

	case core.EnclaveProviderSpotify:
		html, err := object.GetSpotifyWidgetHtml(enc)
		if err != nil || html == "" {
			// html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object spotify-enclave-object error">Failed to load spotify widget from %s</div></div>`, enc.ObjectID)
			html = r.wrapEnclaveErrorHtml("spotify", enc.ObjectID)
		} else {
			// html = fmt.Sprintf(`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object spotify-enclave-object normal-object no-border">%s</div></div>`, html)
			html = r.wrapEnclaveHtml("spotify", html, true, false)
		}
		w.Write([]byte(html))

	case core.EnclaveHtml5Audio:
		html, err := object.GetAudioHtml(enc)
		if err != nil || html == "" {
			html = r.wrapEnclaveErrorHtml("audio", enc.ObjectID)
		} else {
			html = r.wrapEnclaveHtml("audio", html, true, false)
		}
		w.Write([]byte(html))

	case core.EnclaveProviderQuailImage:
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
			html = r.wrapEnclaveErrorHtml("quail-image", enc.ObjectID)
		}
		w.Write([]byte(html))

	case core.EnclaveRegularImage:
		alt := string(node.Text(source))
		if alt == "" && len(enc.Title) != 0 {
			alt = fmt.Sprintf("An image to describe %s", enc.Title)
		}
		if alt == "" {
			alt = "An image to describe post"
		}
		html := fmt.Sprintf(`<img src="%s" alt="%s" />`, enc.URL.String(), alt)
		w.Write([]byte(html))

	}

	return ast.WalkContinue, nil
}

func (r *HTMLRenderer) wrapEnclaveErrorHtml(enclaveName, objectID string) string {
	html := fmt.Sprintf(
		`<div class="enclave-object-wrapper normal-wrapper"><div class="enclave-object %s-enclave-object error">Failed to load %s from %s</div></div>`,
		enclaveName, enclaveName, objectID,
	)
	return html
}

func (r *HTMLRenderer) wrapEnclaveHtml(enclaveName, html string, isNormal, hasBorder bool) string {
	normalCls := ""
	borderCls := ""
	autoResizeCls := "normal-wrapper"
	if isNormal {
		normalCls = "normal-object"
	} else {
		autoResizeCls = "auto-resize"
	}
	if !hasBorder {
		borderCls = "no-border"
	}

	ret := fmt.Sprintf(
		`<div class="enclave-object-wrapper %s"><div class="enclave-object %s-enclave-object %s %s">%s</div></div>`,
		autoResizeCls, enclaveName, normalCls, borderCls, html,
	)
	return ret
}
