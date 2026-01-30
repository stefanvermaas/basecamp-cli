package commands

import "testing"

func TestStripHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "no html",
			input: "plain text",
			want:  "plain text",
		},
		{
			name:  "simple tag",
			input: "<p>paragraph</p>",
			want:  "paragraph",
		},
		{
			name:  "nested tags",
			input: "<div><p>nested <strong>content</strong></p></div>",
			want:  "nested content",
		},
		{
			name:  "multiple whitespace",
			input: "<p>first</p>   <p>second</p>",
			want:  "first second",
		},
		{
			name:  "self closing tags",
			input: "line1<br/>line2",
			want:  "line1 line2",
		},
		{
			name:  "attributes",
			input: `<a href="http://example.com">link</a>`,
			want:  "link",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripHTML(tt.input)
			if got != tt.want {
				t.Errorf("stripHTML(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCoalesce(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   string
	}{
		{
			name:   "all empty",
			values: []string{"", "", ""},
			want:   "",
		},
		{
			name:   "first non-empty",
			values: []string{"first", "second", "third"},
			want:   "first",
		},
		{
			name:   "second non-empty",
			values: []string{"", "second", "third"},
			want:   "second",
		},
		{
			name:   "last non-empty",
			values: []string{"", "", "third"},
			want:   "third",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := coalesce(tt.values...)
			if got != tt.want {
				t.Errorf("coalesce(%v) = %q, want %q", tt.values, got, tt.want)
			}
		})
	}
}
