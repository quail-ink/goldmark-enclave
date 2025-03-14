package callout

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
)

// A Callout struct represents a GitHub-style callout block.
type Callout struct {
	ast.BaseBlock
	Title   string
	content bytes.Buffer
}

// Dump implements Node.Dump.
func (n *Callout) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindCallout is a NodeKind of the Callout node.
var KindCallout = ast.NewNodeKind("Callout")

// Kind implements Node.Kind.
func (n *Callout) Kind() ast.NodeKind {
	return KindCallout
}

// NewCallout returns a new Callout node.
func NewCallout() *Callout {
	return &Callout{
		BaseBlock: ast.BaseBlock{},
	}
}

// SetTitle sets the title of the callout.
func (n *Callout) SetTitle(title string) {
	n.Title = title
}

// AppendContent appends content to the callout.
func (n *Callout) AppendContent(content []byte) {
	n.content.Write(content)
}

// Content returns the content of the callout.
func (n *Callout) Content() []byte {
	return n.content.Bytes()
}
