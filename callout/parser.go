package callout

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type calloutParser struct {
}

var defaultCalloutParser = &calloutParser{}

// NewCalloutParser returns a new BlockParser for callouts.
func NewCalloutParser() parser.BlockParser {
	return defaultCalloutParser
}

var calloutRegexp = regexp.MustCompile(`^\>\s*\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]\s*$`)

func (b *calloutParser) Trigger() []byte {
	return []byte{'>'}
}

func (b *calloutParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	node := NewCallout()

	line, _ := reader.PeekLine()
	if !calloutRegexp.Match(line) {
		return nil, parser.NoChildren
	}

	matches := calloutRegexp.FindStringSubmatch(string(line))
	calloutType := matches[1]
	reader.Advance(len(line) - 1)

	// Map callout types
	title := calloutType
	blockType := strings.ToLower(calloutType)
	if calloutType == "NOTE" {
		blockType = "info"
	} else if calloutType == "CAUTION" {
		blockType = "danger"
	}

	node.SetTitle(title)
	node.SetAttribute([]byte("class"), []byte(fmt.Sprintf("custom-block %s", blockType)))
	node.SetAttribute([]byte("data-title"), []byte(title))
	node.SetAttribute([]byte("data-type"), []byte(blockType))

	return node, parser.HasChildren
}

func (b *calloutParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	// Check if line starts with >
	if len(line) > 0 && line[0] == '>' {
		// This is still part of the callout
		reader.Advance(1)
	} else {
		return parser.Close
	}

	return parser.Continue | parser.HasChildren
}

func (b *calloutParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
}

func (b *calloutParser) CanInterruptParagraph() bool {
	return true
}

func (b *calloutParser) CanAcceptIndentedLine() bool {
	return false
}
