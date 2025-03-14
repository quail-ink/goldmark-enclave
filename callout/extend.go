package callout

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Extender allows you to use GitHub-style callouts in markdown
//
// Callouts are blockquotes with special syntax:
//
// > [!NOTE]
// > This is a note
//
// > [!TIP]
// > This is a tip
//
// > [!IMPORTANT]
// > This is important
//
// > [!WARNING]
// > This is a warning
//
// > [!CAUTION]
// > This is a caution
type Extender struct {
	priority int // optional int != 0. the priority value for parser and renderer. Defaults to 100.
}

func New() *Extender {
	return &Extender{
	}
}

// This implements the Extend method for goldmark-callouts.Extender
func (e *Extender) Extend(md goldmark.Markdown) {
	priority := 100

	if e.priority != 0 {
		priority = e.priority
	}
	md.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(&calloutParser{}, priority),
		),
	)
	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&Renderer{}, priority),
		),
	)
}
