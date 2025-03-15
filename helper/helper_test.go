package helper

import (
	"reflect"
	"testing"
)

func TestGetParagraphs(t *testing.T) {
	tests := []struct {
		name     string
		markdown []byte
		want     []string
	}{
		{
			name:     "single paragraph",
			markdown: []byte("This is a simple paragraph."),
			want:     []string{"This is a simple paragraph."},
		},
		{
			name:     "multiple paragraphs",
			markdown: []byte("First paragraph.\n\nSecond paragraph.\n\nThird paragraph."),
			want:     []string{"First paragraph.", "Second paragraph.", "Third paragraph."},
		},
		{
			name:     "skip blockquote",
			markdown: []byte("Normal paragraph.\n\n> This is in a blockquote\n> Another line in blockquote\n\nAnother normal paragraph."),
			want:     []string{"Normal paragraph.", "Another normal paragraph."},
		},
		{
			name:     "skip code blocks",
			markdown: []byte("Before code.\n\n```\ncode block content\nmore code\n```\n\nAfter code."),
			want:     []string{"Before code.", "After code."},
		},
		{
			name:     "skip tables",
			markdown: []byte("Before table.\n\n| Header1 | Header2 |\n|---------|----------|\n| Cell1   | Cell2   |\n\nAfter table."),
			want:     []string{"Before table.", "After table."},
		},
		{
			name:     "empty input",
			markdown: []byte(""),
			want:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetParagraphs(tt.markdown)
			if len(got) != len(tt.want) {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetParagraphs() = %v, %s, want %v, %s, name: %s", got, reflect.TypeOf(got), tt.want, reflect.TypeOf(tt.want), tt.name)
				}
			}
		})
	}
}
