package enclave

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Option func(*enclaveExtension)

type enclaveExtension struct{}

const (
	EnclaveProviderYouTube     = "youtube"
	EnclaveProviderBilibili    = "bilibili"
	EnclaveProviderTwitter     = "twitter"
	EnclaveProviderTradingView = "tradingview"
	EnclaveProviderQuail       = "quail"
)

func New(opts ...Option) goldmark.Extender {
	e := &enclaveExtension{}
	for _, opt := range opts {
		opt(e)
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
			util.Prioritized(NewHTMLRenderer(), 500),
		),
	)
}

type Enclave struct {
	ast.Image
	Provider string
	ObjectID string
	Theme    string
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
