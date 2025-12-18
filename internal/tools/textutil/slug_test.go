package textutil

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic cases
		{
			name:     "simple lowercase",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "simple uppercase",
			input:    "HELLO",
			expected: "hello",
		},
		{
			name:     "mixed case with space",
			input:    "Hello World",
			expected: "hello-world",
		},

		// Special characters
		{
			name:     "punctuation",
			input:    "Hello, World!",
			expected: "hello-world",
		},
		{
			name:     "multiple special chars",
			input:    "Hello!!! World???",
			expected: "hello-world",
		},
		{
			name:     "leading special chars",
			input:    "!!!Hello",
			expected: "hello",
		},
		{
			name:     "trailing special chars",
			input:    "Hello!!!",
			expected: "hello",
		},
		{
			name:     "leading and trailing",
			input:    "---Hello World---",
			expected: "hello-world",
		},

		// Edge cases
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "     ",
			expected: "",
		},
		{
			name:     "numbers",
			input:    "Product 123",
			expected: "product-123",
		},
		{
			name:     "alphanumeric mix",
			input:    "ABC-123-XYZ",
			expected: "abc-123-xyz",
		},

		// Unicode and international characters
		{
			name:     "unicode characters",
			input:    "Café au lait",
			expected: "caf-au-lait",
		},
		{
			name:     "emoji",
			input:    "Hello 👋 World",
			expected: "hello-world",
		},
		{
			name:     "chinese characters",
			input:    "你好 World",
			expected: "world",
		},

		// Real-world examples
		{
			name:     "blog post title",
			input:    "10 Tips for Writing Clean Code!",
			expected: "10-tips-for-writing-clean-code",
		},
		{
			name:     "url with slashes",
			input:    "path/to/resource",
			expected: "path-to-resource",
		},
		{
			name:     "multiple consecutive spaces",
			input:    "Hello    World",
			expected: "hello-world",
		},
		{
			name:     "multiple consecutive spaces",
			input:    "Hello    World",
			expected: "hello-world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Slugify(tt.input)
			if result != tt.expected {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
