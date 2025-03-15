package callout

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/quailyquaily/goldmark-enclave/helper"
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
			expected: fmt.Sprintf(`<div class="custom-block info" data-title="NOTE" data-type="info" data-callout-type="github-style">
<div class="custom-block-title">%sNOTE</div>
<p>This is a note callout.</p>
</div>
`, helper.GetBlockIcon("info")),
		},
		{
			name:     "lowcase tip callout",
			markdown: "> [!tip]\n> This is a lowcase tip callout.",
			expected: fmt.Sprintf(`<div class="custom-block tip" data-title="TIP" data-type="tip" data-callout-type="github-style">
<div class="custom-block-title">%sTIP</div>
<p>This is a lowcase tip callout.</p>
</div>
`, helper.GetBlockIcon("tip")),
		},
		{
			name:     "fallback to default callout",
			markdown: "> [!what]\n> This is a default callout.",
			expected: fmt.Sprintf(`<div class="custom-block info" data-title="INFO" data-type="info" data-callout-type="github-style">
<div class="custom-block-title">%sINFO</div>
<p>This is a default callout.</p>
</div>
`, helper.GetBlockIcon("note")),
		},
		{
			name:     "callout with custom title",
			markdown: "> [!info] This is a custom title\n> This is a callout with custom title.",
			expected: fmt.Sprintf(`<div class="custom-block info" data-title="This is a custom title" data-type="info" data-callout-type="github-style">
<div class="custom-block-title">%sThis is a custom title</div>
<p>This is a callout with custom title.</p>
</div>
`, helper.GetBlockIcon("info")),
		},
		{
			name:     "multi-line callout",
			markdown: "> [!NOTE]\n> First line\n> Second line",
			expected: fmt.Sprintf(`<div class="custom-block info" data-title="NOTE" data-type="info" data-callout-type="github-style">
<div class="custom-block-title">%sNOTE</div>
<p>First line
Second line</p>
</div>
`, helper.GetBlockIcon("info")),
		},
		{
			name:     "simple inline callout",
			markdown: "WARNING: This is a simple inline callout.",
			expected: fmt.Sprintf(`<div class="custom-block warning" data-title="WARNING" data-type="warning" data-callout-type="simple-inline">
<div class="custom-block-title">%sWARNING</div>
<p>This is a simple inline callout.</p>
</div>
`, helper.GetBlockIcon("warning")),
		},
		{
			name:     "simple callout",
			markdown: "WARNING\nThis is a simple callout.",
			expected: fmt.Sprintf(`<div class="custom-block warning" data-title="WARNING" data-type="warning" data-callout-type="simple">
<div class="custom-block-title">%sWARNING</div>
<p>This is a simple callout.</p>
</div>
`, helper.GetBlockIcon("warning")),
		},
		{
			name:     "simple callout with exclamation mark",
			markdown: "WARNING!!!\nThis is a simple callout with exclamation mark.",
			expected: fmt.Sprintf(`<div class="custom-block warning" data-title="WARNING" data-type="warning" data-callout-type="simple">
<div class="custom-block-title">%sWARNING</div>
<p>This is a simple callout with exclamation mark.</p>
</div>
`, helper.GetBlockIcon("warning")),
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
