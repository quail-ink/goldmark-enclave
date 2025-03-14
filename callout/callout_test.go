package callout

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
)

func TestCallouts(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "note callout",
			markdown: "> [!NOTE]\n> This is a note callout.",
			expected: "<div class=\"custom-block info\" data-title=\"NOTE\">\n<div class=\"custom-block-title\">NOTE</div>\n<p>This is a note callout.</p>\n</div>\n",
		},
		{
			name:     "tip callout",
			markdown: "> [!TIP]\n> This is a tip callout.",
			expected: "<div class=\"custom-block tip\" data-title=\"TIP\">\n<div class=\"custom-block-title\">TIP</div>\n<p>This is a tip callout.</p>\n</div>\n",
		},
		{
			name:     "important callout",
			markdown: "> [!IMPORTANT]\n> This is an important callout.",
			expected: "<div class=\"custom-block important\" data-title=\"IMPORTANT\">\n<div class=\"custom-block-title\">IMPORTANT</div>\n<p>This is an important callout.</p>\n</div>\n",
		},
		{
			name:     "warning callout",
			markdown: "> [!WARNING]\n> This is a warning callout.",
			expected: "<div class=\"custom-block warning\" data-title=\"WARNING\">\n<div class=\"custom-block-title\">WARNING</div>\n<p>This is a warning callout.</p>\n</div>\n",
		},
		{
			name:     "caution callout",
			markdown: "> [!CAUTION]\n> This is a caution callout.",
			expected: "<div class=\"custom-block danger\" data-title=\"CAUTION\">\n<div class=\"custom-block-title\">CAUTION</div>\n<p>This is a caution callout.</p>\n</div>\n",
		},
		{
			name:     "multi-line callout",
			markdown: "> [!NOTE]\n> First line\n> Second line",
			expected: "<div class=\"custom-block info\" data-title=\"NOTE\">\n<div class=\"custom-block-title\">NOTE</div>\n<p>First line\nSecond line</p>\n</div>\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			markdown := goldmark.New(
				goldmark.WithExtensions(
					New(),
				),
			)

			var buf bytes.Buffer
			err := markdown.Convert([]byte(test.markdown), &buf)
			if err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			if got := buf.String(); got != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, got)
			}
		})
	}
}
