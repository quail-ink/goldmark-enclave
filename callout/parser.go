package callout

import (
	"bytes"
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

var calloutRegexp = regexp.MustCompile(`^\>\s*\[!(.*)?\]\s*?(.*)\s*$`)
var simpleCalloutRegexp = regexp.MustCompile(`^(NOTE|TIP|IMPORTANT|WARNING|CAUTION)(\s*[:!]+)?\s*(.*)$`)
var supportedCalloutTypeTitle = []string{"NOTE", "INFO", "TIP", "IMPORTANT", "WARNING", "CAUTION"}

func (b *calloutParser) Trigger() []byte {
	// Detect both '>' for GitHub-style and the first letters of supported callout types
	return []byte{'>', 'N', 'T', 'I', 'W', 'C'}
}

func (b *calloutParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	node := NewCallout()

	_line, _ := reader.PeekLine()
	line := bytes.TrimSpace(_line)

	// Check for GitHub-style callout
	if line[0] == '>' && calloutRegexp.Match(line) {
		matches := calloutRegexp.FindStringSubmatch(string(line))
		calloutType := matches[1]
		customTitle := ""
		if len(matches) > 2 {
			customTitle = strings.TrimSpace(matches[2])
		}

		// Map callout types
		blockType := mapCalloutType(calloutType)

		title := mapCalloutTitle(calloutType, customTitle)

		node.SetTitle(title)
		node.SetAttribute([]byte("class"), []byte(fmt.Sprintf("custom-block %s", blockType)))
		node.SetAttribute([]byte("data-title"), []byte(title))
		node.SetAttribute([]byte("data-type"), []byte(blockType))
		node.SetAttribute([]byte("data-callout-type"), []byte("github-style"))

		reader.Advance(len(line))

		return node, parser.HasChildren
	}

	// Check for simple callout syntax
	if simpleCalloutRegexp.Match(line) {
		matches := simpleCalloutRegexp.FindStringSubmatch(string(line))
		calloutType := strings.TrimSpace(matches[1])
		content := strings.TrimSpace(matches[3])

		// Map callout types
		blockType := mapCalloutType(calloutType)
		title := mapCalloutTitle(calloutType, "")

		node.SetTitle(title)
		node.SetAttribute([]byte("class"), []byte(fmt.Sprintf("custom-block %s", blockType)))
		node.SetAttribute([]byte("data-title"), []byte(title))
		node.SetAttribute([]byte("data-type"), []byte(blockType))
		node.SetAttribute([]byte("data-callout-type"), []byte("simple"))

		// If content exists on the same line, add it as a paragraph
		if content != "" {
			node.SetAttribute([]byte("data-callout-type"), []byte("simple-inline"))
			node.SetContent([]byte(content))
			reader.Advance(len(calloutType) + 1)
		} else {
			reader.Advance(len(line))
		}

		return node, parser.HasChildren
	}

	return nil, parser.HasChildren
}

func (b *calloutParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	_line, _ := reader.PeekLine()
	line := bytes.TrimSpace(_line)
	// For GitHub style, check if line starts with '>'
	if callout, ok := node.(*Callout); ok {
		if val, ok2 := callout.Attribute([]byte("data-callout-type")); ok2 {
			calloutType := string(val.([]byte))
			if calloutType == "github-style" {
				if len(line) > 0 && line[0] == '>' {
					// This is still part of the callout
					reader.Advance(1)
					return parser.Continue | parser.HasChildren
				} else {
					return parser.Close
				}
			} else if calloutType == "simple" {
				if len(line) == 0 {
					return parser.Close
				}
			} else if calloutType == "simple-inline" {
				if len(line) == 0 {
					return parser.Close
				}
			}
		}
	}

	// Check if the content is part of a simple callout
	if len(line) == 0 {
		return parser.Close
	}

	// For simple callout, continue until an empty line or another callout
	if simpleCalloutRegexp.Match(line) {
		return parser.Close
	}

	return parser.Continue | parser.HasChildren
}

func (b *calloutParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// Handle any inline content that was saved during Open
	if callout, ok := node.(*Callout); ok {
		content := callout.Content()
		if len(content) > 0 {
			// Content will be processed later as part of the callout
			callout.SetContent([]byte(""))
		}
	}
}

func (b *calloutParser) CanInterruptParagraph() bool {
	return true
}

func (b *calloutParser) CanAcceptIndentedLine() bool {
	return false
}

// Helper function to map callout types to their CSS classes
func mapCalloutType(calloutType string) string {
	calloutType = strings.ToUpper(calloutType)
	blockType := "info"
	if calloutType == "NOTE" {
		blockType = "info"
	} else if calloutType == "CAUTION" {
		blockType = "danger"
	} else if calloutType == "TIP" {
		blockType = "tip"
	} else if calloutType == "IMPORTANT" {
		blockType = "important"
	} else if calloutType == "WARNING" {
		blockType = "warning"
	} else {
		blockType = "info"
	}
	return blockType
}

func mapCalloutTitle(calloutType string, customTitle string) string {
	calloutType = strings.ToUpper(calloutType)
	if customTitle != "" {
		return customTitle
	}
	title := "INFO"
	for _, supportedType := range supportedCalloutTypeTitle {
		if calloutType == supportedType {
			title = supportedType
			break
		}
	}
	return title
}
