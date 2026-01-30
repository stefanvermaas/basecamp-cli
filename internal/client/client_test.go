package client

import (
	"testing"
)

func TestParseNextLink(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{
			name:   "empty",
			header: "",
			want:   "",
		},
		{
			name:   "single next link",
			header: `<https://api.example.com/page2>; rel="next"`,
			want:   "https://api.example.com/page2",
		},
		{
			name:   "multiple links with next",
			header: `<https://api.example.com/page1>; rel="prev", <https://api.example.com/page3>; rel="next"`,
			want:   "https://api.example.com/page3",
		},
		{
			name:   "only prev link",
			header: `<https://api.example.com/page1>; rel="prev"`,
			want:   "",
		},
		{
			name:   "complex URL",
			header: `<https://3.basecampapi.com/12345/buckets/67890/comments.json?page=2>; rel="next"`,
			want:   "https://3.basecampapi.com/12345/buckets/67890/comments.json?page=2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseNextLink(tt.header)
			if got != tt.want {
				t.Errorf("parseNextLink(%q) = %q, want %q", tt.header, got, tt.want)
			}
		})
	}
}

func TestResolveURL(t *testing.T) {
	c := &Client{
		baseURL: "https://3.basecampapi.com/12345",
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "relative path",
			path: "/projects.json",
			want: "https://3.basecampapi.com/12345/projects.json",
		},
		{
			name: "absolute URL",
			path: "https://other.api.com/data.json",
			want: "https://other.api.com/data.json",
		},
		{
			name: "http URL",
			path: "http://localhost:3000/test",
			want: "http://localhost:3000/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.resolveURL(tt.path)
			if got != tt.want {
				t.Errorf("resolveURL(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
