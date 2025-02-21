package href

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/quail-ink/goldmark-enclave/helper"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type (
	HrefExtension struct {
		cfg *Config
	}
	hrefRenderer struct {
		cfg *Config
	}
	Config struct {
		NoFollowByDefault bool
		DoFollowDomains   []string
	}
)

func New(cfg *Config) *HrefExtension {
	return &HrefExtension{
		cfg: cfg,
	}
}

func (e *HrefExtension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewHrefRenderer(e.cfg), 500),
	))
}

func NewHrefRenderer(cfg *Config) renderer.NodeRenderer {
	return &hrefRenderer{cfg: cfg}
}

func (r *hrefRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderHref)
	reg.Register(ast.KindAutoLink, r.renderHref)
}

func (r *hrefRenderer) renderHref(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		_, _ = w.Write([]byte("</a>"))
		return ast.WalkContinue, nil
	}

	title := ""
	dst := ""

	isLink := false
	switch n := node.(type) {
	case *ast.Link:
		{
			dst = string(n.Destination)
			title = string(n.Title)
			if title == "" {
				anchorText := helper.ExtractTextFromNode(n, source)
				if anchorText != "" {
					title = string(anchorText)
				} else {
					u, err := url.Parse(dst)
					if err == nil {
						title = fmt.Sprintf("A Link to %s", u.Host)
					} else {
						title = "A Link"
					}
				}
			}
			isLink = true
		}
	case *ast.AutoLink:
		dst = string(n.URL(source))
		title = fmt.Sprintf("A Link of %s", dst)
		isLink = false
	default:
		return ast.WalkContinue, nil
	}

	rel := []string{"rel", "noopener"}
	if r.cfg.NoFollowByDefault {
		rel = []string{"rel", "noopener ugc nofollow"}
		if r.shouldFollowLink(dst) {
			rel = []string{"rel", "noopener"}
		}
	}
	attrs := [][]string{
		{"href", dst},
		{"title", title},
		rel,
	}
	tag := helper.HTMLTag("a", attrs)

	_, _ = w.Write(tag)

	if !isLink {
		_, _ = w.Write([]byte(dst))
	}
	return ast.WalkContinue, nil
}

func (r *hrefRenderer) shouldFollowLink(link string) bool {
	if len(r.cfg.DoFollowDomains) == 0 {
		return false
	}
	u, err := url.Parse(link)
	if err != nil {
		return false
	}
	for _, domain := range r.cfg.DoFollowDomains {
		if strings.EqualFold(u.Host, domain) {
			return true
		}
	}
	return false
}
