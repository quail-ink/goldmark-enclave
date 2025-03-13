package kbd

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"
)

func TestKbd(t *testing.T) {
	// Create Goldmark with our extension
	markdown := goldmark.New(
		goldmark.WithExtensions(
			New(),
		),
		// No need for WithUnsafe() now
	)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple kbd",
			input:    "<kbd>Ctrl</kbd>",
			expected: "<p><kbd>Ctrl</kbd></p>\n",
		},
		{
			name:     "Kbd with plus sign",
			input:    "<kbd>Ctrl</kbd> + <kbd>A</kbd>",
			expected: "<p><kbd>Ctrl</kbd> + <kbd>A</kbd></p>\n",
		},
		{
			name:     "Multiple kbd tags",
			input:    "<kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>Esc</kbd>",
			expected: "<p><kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>Esc</kbd></p>\n",
		},
		{
			name:     "Kbd in a sentence",
			input:    "Press <kbd>Enter</kbd> to continue.",
			expected: "<p>Press <kbd>Enter</kbd> to continue.</p>\n",
		},
		{
			name:     "Other HTML tags should be sanitized",
			input:    "<div>This should be sanitized</div> \nbut <kbd>Ctrl</kbd> should remain",
			expected: "<p><!-- raw HTML omitted -->This should be sanitized<!-- raw HTML omitted --> but <kbd>Ctrl</kbd> should remain</p>\n",
		},
		{
			name:     "Incomplete kbd tag",
			input:    "<kbd>Ctrl",
			expected: "<p><kbd>Ctrl</p>\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := markdown.Convert([]byte(test.input), &buf); err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			if buf.String() != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, buf.String())
			}
		})
	}
}

func TestKbdSource(t *testing.T) {
	// This test verifies that the extension works with the example source files
	testutil.DoTestCaseFile(goldmark.New(
		goldmark.WithExtensions(
			New(),
		),
		// No need to disable HTML sanitization
	), "_testdata/kbd.txt", t)
}
