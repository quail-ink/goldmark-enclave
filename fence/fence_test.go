package fence

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestFencedContainers(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name: "basic fenced container",
			markdown: `:::info
This is a basic info container.
:::`,
			expected: "<div data-fence=\"0\" class=\"custom-block info\" data-title=\"INFO\">\n<div class=\"custom-block-title\">INFO</div>\n<p>This is a basic info container.</p>\n</div>\n",
		},
		{
			name: "container with custom title",
			markdown: `:::info Custom Title
This container has a custom title.
:::`,
			expected: "<div data-fence=\"0\" class=\"custom-block info\" data-title=\"Custom Title\">\n<div class=\"custom-block-title\">Custom Title</div>\n<p>This container has a custom title.</p>\n</div>\n",
		},
		{
			name: "warning container",
			markdown: `:::warning
This is a warning container.
:::`,
			expected: "<div data-fence=\"0\" class=\"custom-block warning\" data-title=\"WARNING\">\n<div class=\"custom-block-title\">WARNING</div>\n<p>This is a warning container.</p>\n</div>\n",
		},
		{
			name: "nested containers",
			markdown: `:::info
Outer container

:::warning
Inner warning container
:::

More outer container content
:::`,
			expected: "<div data-fence=\"0\" class=\"custom-block info\" data-title=\"INFO\">\n<div class=\"custom-block-title\">INFO</div>\n<p>Outer container</p>\n<div data-fence=\"1\" class=\"custom-block warning\" data-title=\"WARNING\">\n<div class=\"custom-block-title\">WARNING</div>\n<p>Inner warning container</p>\n</div>\n<p>More outer container content</p>\n</div>\n",
		},
		{
			name: "container with multiple block elements",
			markdown: `:::tip
# Heading inside container

- List item 1
- List item 2

> Blockquote inside container

` + "```go" + `
func example() {
    fmt.Println("Code block inside container")
}
` + "```" + `
:::`,
			expected: "<div data-fence=\"0\" class=\"custom-block tip\" data-title=\"TIP\">\n<div class=\"custom-block-title\">TIP</div>\n<h1>Heading inside container</h1>\n<ul>\n<li>List item 1</li>\n<li>List item 2</li>\n</ul>\n<blockquote>\n<p>Blockquote inside container</p>\n</blockquote>\n<pre><code class=\"language-go\">func example() {\n    fmt.Println(&quot;Code block inside container&quot;)\n}\n</code></pre>\n</div>\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			markdown := goldmark.New(
				goldmark.WithExtensions(
					&Extender{},
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

func TestComplexDocument(t *testing.T) {
	markdown := `# Document with Fenced Containers

Regular paragraph.

:::{.info}
This is an information container.

## Heading inside container

Some more text.
:::

More regular content.

:::{#special-container .warning}
This is a warning with an ID.

:::{.danger}
Nested danger container.
:::

End of warning.
:::

Final paragraph.`

	expected := `<h1>Document with Fenced Containers</h1>
<p>Regular paragraph.</p>
<div class="custom-block info" data-title="INFO">
<div class="custom-block-title">INFO</div>
<p>This is an information container.</p>
<h2>Heading inside container</h2>
<p>Some more text.</p>
</div>
<p>More regular content.</p>
<div id="special-container" class="custom-block warning" data-title="WARNING">
<div class="custom-block-title">WARNING</div>
<p>This is a warning with an ID.</p>
<div class="custom-block danger" data-title="DANGER">
<div class="custom-block-title">DANGER</div>
<p>Nested danger container.</p>
</div>
<p>End of warning.</p>
</div>
<p>Final paragraph.</p>
`

	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			&Extender{},
		),
	)

	var buf bytes.Buffer
	err := md.Convert([]byte(markdown), &buf)
	if err != nil {
		t.Fatalf("Failed to convert markdown: %v", err)
	}

	if got := buf.String(); got != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, got)
	}
}

// Add a benchmark test for fenced container parsing
func BenchmarkFencedContainers(b *testing.B) {
	source := []byte(`# Document with Fenced Containers

:::{.info}
This is an information container.

## Heading inside container

Some more text.
:::

:::{#special-container .warning}
This is a warning with an ID.

:::{.danger}
Nested danger container.
:::

End of warning.
:::`)

	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			&Extender{},
		),
	)

	var buf bytes.Buffer
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := markdown.Convert(source, &buf)
		if err != nil {
			b.Fatalf("Failed to convert markdown: %v", err)
		}
	}
}
