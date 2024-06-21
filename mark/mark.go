package mark

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type (
	MarkExtension    struct{}
	MarkHTMLRenderer struct{}
	MarkParser       struct{}
	MarkASTNode      struct {
		ast.BaseInline
		Content []byte
	}
)

func New() *MarkExtension {
	return &MarkExtension{}
}

func (e *MarkExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewMarkParser(), 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewMarkHTMLRenderer(), 500),
	))
}

var KindMark = ast.NewNodeKind("Mark")

func (n *MarkASTNode) Kind() ast.NodeKind {
	return KindMark
}

func (n *MarkASTNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

func NewMarkParser() parser.InlineParser {
	return &MarkParser{}
}

func (p *MarkParser) Trigger() []byte {
	return []byte{'='}
}

// Parse parses the ==highlight== syntax and adds the MarkASTNode to the AST.
func (p *MarkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()
	pos := segment.Start
	end := pos

	if len(line) < 4 || line[0] != '=' || line[1] != '=' {
		return nil
	}

	for i := 2; i < len(line)-1; i++ {
		if line[i] == '=' && line[i+1] == '=' {
			end = segment.Start + i + 2
			break
		}
	}

	if end == pos {
		return nil
	}

	block.Advance(end - pos)
	content := line[2 : end-pos-2]
	return &MarkASTNode{Content: content}
}

func NewMarkHTMLRenderer() renderer.NodeRenderer {
	return &MarkHTMLRenderer{}
}

func (r *MarkHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMark, r.renderMark)
}

func (r *MarkHTMLRenderer) renderMark(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*MarkASTNode)
	if entering {
		_, _ = w.WriteString("<mark>")
		_, _ = w.Write(n.Content)
	} else {
		_, _ = w.WriteString("</mark>")
	}
	return ast.WalkContinue, nil
}
