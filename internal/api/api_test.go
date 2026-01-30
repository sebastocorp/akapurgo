package api

import (
	"testing"
)

func TestAddQueryParam(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		key      string
		value    string
		expected string
		wantErr  bool
	}{
		{
			name:     "URL without existing query params",
			url:      "https://example.com/page",
			key:      "imbypass",
			value:    "true",
			expected: "https://example.com/page?imbypass=true",
			wantErr:  false,
		},
		{
			name:     "URL with existing query params",
			url:      "https://example.com/page?foo=bar",
			key:      "imbypass",
			value:    "true",
			expected: "https://example.com/page?foo=bar&imbypass=true",
			wantErr:  false,
		},
		{
			name:     "URL with multiple existing query params",
			url:      "https://example.com/page?foo=bar&baz=qux",
			key:      "imbypass",
			value:    "true",
			expected: "https://example.com/page?baz=qux&foo=bar&imbypass=true",
			wantErr:  false,
		},
		{
			name:     "URL with path only",
			url:      "/path/to/resource",
			key:      "imbypass",
			value:    "true",
			expected: "/path/to/resource?imbypass=true",
			wantErr:  false,
		},
		{
			name:     "Invalid URL",
			url:      "://invalid",
			key:      "imbypass",
			value:    "true",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := addQueryParam(tt.url, tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("addQueryParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("addQueryParam() = %v, want %v", result, tt.expected)
			}
		})
	}
}
