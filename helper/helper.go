package helper

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
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

func ExtractTextPartsFromNode(node ast.Node, source []byte) []string {
	parts := make([]string, 0)
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			parts = append(parts, string(textNode.Segment.Value(source)))
		}
	}
	return parts
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

func GetParagraphs(source []byte) []string {
	parser := goldmark.DefaultParser()
	reader := text.NewReader(source)
	root := parser.Parse(reader)

	var paragraphs []string

	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() == ast.KindParagraph {
			parent := n.Parent()
			// Skip paragraphs inside blockquotes or code blocks
			if parent != nil && (parent.Kind() == ast.KindBlockquote || parent.Kind() == ast.KindFencedCodeBlock) {
				return ast.WalkContinue, nil
			}

			lines := n.Lines()
			var paragraphText bytes.Buffer
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				paragraphText.Write(line.Value(source))
			}
			paragraphs = append(paragraphs, paragraphText.String())
		}

		return ast.WalkContinue, nil
	})

	return paragraphs
}

func ConvertKindParagraphsToNormal(source []byte, paras []string) []byte {
	sourceStr := string(source)
	for _, para := range paras {
		c := strings.ReplaceAll(para, "\n", "\n\n")
		sourceStr = strings.ReplaceAll(sourceStr, para, c)
	}
	return []byte(sourceStr)
}
