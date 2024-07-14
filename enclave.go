package enclave

import (
	"net/url"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type (
	Option           func(*enclaveExtension)
	enclaveExtension struct {
		cfg *Config
	}

	Config struct {
		DefaultImageAltPrefix string
	}
)

const (
	EnclaveProviderYouTube     = "youtube"
	EnclaveProviderBilibili    = "bilibili"
	EnclaveProviderTwitter     = "twitter"
	EnclaveProviderTradingView = "tradingview"
	EnclaveProviderDifyWidget          = "dify-widget"
	EnclaveProviderQuailWidget = "quail-widget"
	EnclaveProviderQuailImage  = "quail-image"
	EnclaveRegularImage        = "regular-image"
)

func New(cfg *Config) goldmark.Extender {
	e := &enclaveExtension{
		cfg: cfg,
	}
	return e
}

func (e *enclaveExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(defaultASTTransformer, 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewHTMLRenderer(e.cfg), 500),
		),
	)
}

type Enclave struct {
	ast.Image
	URL      *url.URL
	Provider string
	ObjectID string
	Theme    string
	Params   map[string]string
}

var KindEnclave = ast.NewNodeKind("Enclave")

func (n *Enclave) Kind() ast.NodeKind {
	return KindEnclave
}

func NewEnclave(c *Enclave) *Enclave {
	c.Destination = c.Image.Destination
	c.Title = c.Image.Title
	return c
}
