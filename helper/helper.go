package helper

import (
	"fmt"
	"html"
	"strings"

	"github.com/yuin/goldmark/ast"
)

func HTMLTag(tagName string, attrs [][]string) []byte {
	var sb strings.Builder
	sb.WriteString("<")
	sb.WriteString(tagName)
	for _, attr := range attrs {
		if len(attr) >= 2 {
			sb.WriteString(fmt.Sprintf(` %s="%s"`, attr[0], html.EscapeString(attr[1])))
		}
	}
	sb.WriteString(">")
	return []byte(sb.String())
}

func ExtractTextFromNode(node ast.Node, source []byte) string {
	var text []byte
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			text = append(text, textNode.Segment.Value(source)...)
		}
	}
	return string(text)
}
