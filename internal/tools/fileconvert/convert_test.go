package fileconvert

import (
	"io"
	"strings"
	"testing"
)

func TestToUpperText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple lowercase",
			input:    "hello",
			expected: "HELLO",
		},
		{
			name:     "mixed case",
			input:    "Hello World",
			expected: "HELLO WORLD",
		},
		{
			name:     "already uppercase",
			input:    "HELLO",
			expected: "HELLO",
		},
		{
			name:     "with numbers",
			input:    "hello123",
			expected: "HELLO123",
		},
		{
			name:     "with punctuation",
			input:    "Hello, World!",
			expected: "HELLO, WORLD!",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "multiline",
			input:    "hello\nworld",
			expected: "HELLO\nWORLD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)

			result, err := ToUpperText(reader)
			if err != nil {
				t.Fatalf("ToUpperText failed: %v", err)
			}

			output, err := io.ReadAll(result)
			if err != nil {
				t.Fatalf("Failed to read result: %v", err)
			}

			if string(output) != tt.expected {
				t.Errorf("ToUpperText(%q) = %q, want %q", tt.input, string(output), tt.expected)
			}
		})
	}
}
