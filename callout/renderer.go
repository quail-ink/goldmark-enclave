package callout

import (
	"fmt"

	"github.com/quailyquaily/goldmark-enclave/helper"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// A Config struct has configurations for the HTML based renderers.
type Config struct {
	Writer    html.Writer
	HardWraps bool
	XHTML     bool
	Unsafe    bool
}

// CalloutAttributeFilter defines attribute names which callout elements can have
var CalloutAttributeFilter = html.GlobalAttributeFilter

// A Renderer struct is an implementation of renderer.NodeRenderer that renders
// nodes as (X)HTML.
type Renderer struct {
	Config
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindCallout, r.renderCallout)
}

func (r *Renderer) renderCallout(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*Callout)
	if entering {
		_, _ = w.WriteString("<div")
		html.RenderAttributes(w, n, CalloutAttributeFilter)
		_, _ = w.WriteString(">\n")
		icon := ""
		blockType, ok := n.AttributeString("data-type")
		if ok {
			icon = helper.GetBlockIcon(string(blockType.([]byte)))
		}
		_, _ = w.WriteString(fmt.Sprintf("<div class=\"custom-block-title\">%s%s</div>\n", icon, n.Title))
	} else {
		_, _ = w.WriteString("</div>\n")
	}
	return ast.WalkContinue, nil
}
