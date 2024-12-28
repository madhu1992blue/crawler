package main

import (
	"reflect"
	"testing"
)

func TestURLParserURL(t *testing.T) {
	tests := []struct {
		name       string
		htmlBody   string
		rawBaseURL string
		expected   []string
	}{
		{
			name:       "Simple link",
			htmlBody:   `<a href="https://google.com"></a>`,
			rawBaseURL: "https://example.com",
			expected:   []string{"https://google.com"},
		},
		{
			name:       "Nested Link",
			htmlBody:   `<body><a href="https://google.com"></a></body>`,
			rawBaseURL: "https://example.com",
			expected:   []string{"https://google.com"},
		},
		{
			name:       "Multiple Link",
			htmlBody:   `<body><a href="https://example.com/link1"></a><a href="https://example.com/link2"></a></body>`,
			rawBaseURL: "https://example.com",
			expected:   []string{"https://example.com/link1", "https://example.com/link2"},
		},
		{
			name:       "Relative Link",
			htmlBody:   `<body><a href="/link1"></a><a href="/link2"></a></body>`,
			rawBaseURL: "https://example.com",
			expected:   []string{"https://example.com/link1", "https://example.com/link2"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.htmlBody, tc.rawBaseURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
