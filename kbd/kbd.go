package kbd

import (
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type KbdNode struct {
	ast.BaseInline
}

// Dump implements Node.Dump interface.
func (n *KbdNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// Kind implements Node.Kind interface.
func (n *KbdNode) Kind() ast.NodeKind {
	return NodeKind
}

// NodeKind is a unique identifier for the KbdNode.
var NodeKind = ast.NewNodeKind("Kbd")

// New returns a new Extension that enables <kbd> tag parsing and rendering.
func New() goldmark.Extender {
	return &kbdExtension{}
}

type kbdExtension struct{}

// Extend implements goldmark.Extender interface.
func (e *kbdExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewParser(), 100),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewRenderer(), 500),
			util.Prioritized(NewKbdHTMLRenderer(), 501),
		),
	)
}

// KbdParser is a parser for <kbd> tags.
type KbdParser struct{}

// NewParser returns a new Parser for <kbd> tags.
func NewParser() parser.InlineParser {
	return &KbdParser{}
}

// Trigger returns character that triggers this parser.
func (p *KbdParser) Trigger() []byte {
	return []byte{'<'}
}

// Parse parses the <kbd> tag syntax.
func (p *KbdParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()
	if len(line) < 5 { // Minimum for <kbd>
		return nil
	}

	if string(line[:5]) != "<kbd>" {
		return nil
	}

	// Find the closing tag
	closingPos := -1
	for i := 5; i < len(line)-6; i++ {
		if string(line[i:i+6]) == "</kbd>" {
			closingPos = i
			break
		}
	}

	if closingPos == -1 {
		return nil
	}

	// Consume the opening tag
	block.Advance(5)

	// Create the node
	node := &KbdNode{}

	// Parse the content inside the kbd tag
	contentSegment := text.NewSegment(segment.Start+5, segment.Start+closingPos)
	content := block.Value(contentSegment)
	textNode := ast.NewString(content)
	node.AppendChild(node, textNode)

	// Consume the content and closing tag
	block.Advance(closingPos - 5 + 6)

	return node
}

// KbdRenderer is a renderer for KbdNode.
type KbdRenderer struct {
	html.Config
}

// NewRenderer returns a new KbdRenderer.
func NewRenderer() renderer.NodeRenderer {
	return &KbdRenderer{
		Config: html.NewConfig(),
	}
}

// RegisterFuncs registers renderer functions for KbdNode.
func (r *KbdRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(NodeKind, r.renderKbd)
}

func (r *KbdRenderer) renderKbd(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<kbd>")
	} else {
		_, _ = w.WriteString("</kbd>")
	}
	return ast.WalkContinue, nil
}

// KbdHTMLRenderer is a custom HTML renderer that selectively allows <kbd> tags
type KbdHTMLRenderer struct {
	html.Config
}

// NewKbdHTMLRenderer returns a new HTML renderer that allows <kbd> tags
func NewKbdHTMLRenderer() renderer.NodeRenderer {
	return &KbdHTMLRenderer{
		Config: html.NewConfig(),
	}
}

// RegisterFuncs registers rendering functions for HTML nodes
func (r *KbdHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
}

var kbdTagRegexp = regexp.MustCompile(`</?kbd>`)

// renderRawHTML renders raw HTML nodes, allowing only <kbd> tags
func (r *KbdHTMLRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.RawHTML)
	segments := n.Segments
	for i := 0; i < segments.Len(); i++ {
		segment := segments.At(i)
		rawHTML := segment.Value(source)

		// Check if this is a <kbd> tag
		if kbdTagRegexp.Match(rawHTML) {
			_, _ = w.Write(rawHTML)
		} else {
			_, _ = w.Write(util.EscapeHTML(rawHTML))
		}
	}

	return ast.WalkContinue, nil
}
