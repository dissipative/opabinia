package parser

import "testing"

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "MarkdownWithLevelOneHeading",
			input:    []byte("# Title Here\nSome content here."),
			expected: "Title Here",
		},
		{
			name:     "MarkdownWithMultipleHeadings",
			input:    []byte("# First Title\n## Subheading\n# Another Title"),
			expected: "First Title", // Expect the first level 1 heading
		},
		{
			name:     "MarkdownWithoutLevelOneHeading",
			input:    []byte("Regular content without a title"),
			expected: "",
		},
		{
			name:     "EmptyMarkdown",
			input:    []byte(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTitle(tt.input)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
