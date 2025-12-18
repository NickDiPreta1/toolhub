package encodingutil

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text",
			input:    "hello",
			expected: "aGVsbG8=",
		},
		{
			name:     "text with spaces",
			input:    "hello world",
			expected: "aGVsbG8gd29ybGQ=",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},

		// Special characters
		{
			name:     "special characters",
			input:    "Hello, World!",
			expected: "SGVsbG8sIFdvcmxkIQ==",
		},
		{
			name:     "punctuation",
			input:    "foo@bar.com",
			expected: "Zm9vQGJhci5jb20=",
		},

		// Numbers
		{
			name:     "numbers",
			input:    "12345",
			expected: "MTIzNDU=",
		},
		{
			name:     "alphanumeric",
			input:    "abc123",
			expected: "YWJjMTIz",
		},

		// Unicode and multi-byte
		{
			name:     "unicode emoji",
			input:    "Hello 👋",
			expected: "SGVsbG8g8J+Riw==",
		},
		{
			name:     "unicode accents",
			input:    "Café",
			expected: "Q2Fmw6k=",
		},

		// Longer strings
		{
			name:     "sentence",
			input:    "The quick brown fox jumps over the lazy dog",
			expected: "VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZw==",
		},

		// Edge cases
		{
			name:     "single character",
			input:    "a",
			expected: "YQ==",
		},
		{
			name:     "newline",
			input:    "hello\nworld",
			expected: "aGVsbG8Kd29ybGQ=",
		},
		{
			name:     "tab",
			input:    "hello\tworld",
			expected: "aGVsbG8Jd29ybGQ=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.input)
			if got != tt.expected {
				t.Errorf("Encode (%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		// Basic cases (inverse of Encode tests)
		{
			name:        "simple text",
			input:       "aGVsbG8=",
			expected:    "hello",
			expectError: false,
		},
		{
			name:        "text with spaces",
			input:       "aGVsbG8gd29ybGQ=",
			expected:    "hello world",
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: false,
		},

		// Valid Base64
		{
			name:        "special characters",
			input:       "SGVsbG8sIFdvcmxkIQ==",
			expected:    "Hello, World!",
			expectError: false,
		},
		{
			name:        "email",
			input:       "Zm9vQGJhci5jb20=",
			expected:    "foo@bar.com",
			expectError: false,
		},
		{
			name:        "numbers",
			input:       "MTIzNDU=",
			expected:    "12345",
			expectError: false,
		},

		// Unicode
		{
			name:        "unicode emoji",
			input:       "SGVsbG8g8J+Riw==",
			expected:    "Hello 👋",
			expectError: false,
		},
		{
			name:        "unicode accents",
			input:       "Q2Fmw6k=",
			expected:    "Café",
			expectError: false,
		},

		// Error cases - Invalid Base64
		{
			name:        "invalid characters",
			input:       "hello!!!",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid length (not padded)",
			input:       "aGVsbG8", // Missing padding
			expected:    "",
			expectError: true,
		},
		{
			name:        "plain text (not base64)",
			input:       "this is not base64",
			expected:    "",
			expectError: true,
		},
		{
			name:        "special chars not in base64 alphabet",
			input:       "abc@def#",
			expected:    "",
			expectError: true,
		},
		{
			name:        "whitespace in middle",
			input:       "aGVs bG8=",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Decode(%q) expected error but got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Decode(%q) unexpected error: %v", tt.input, err)
			}

			if got != tt.expected {
				t.Errorf("Decode(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestEncodeDecodeRountTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", "hello"},
		{"spaces", "hello world"},
		{"special chars", "Hello, World! #123"},
		{"unicode", "Café ☕"},
		{"emoji", "Hello 👋 World 🌍"},
		{"long text", "The quick brown fox jumps over the lazy dog. Pack my box with five dozen liquor jugs."},
		{"multiline", "line1\nline2\nline3"},
		{"numbers", "0123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Encode(tt.input)
			decoded, err := Decode(encoded)
			if err != nil {
				t.Fatalf("Decode failed after Encode: %v", err)
			}

			if decoded != tt.input {
				t.Errorf("Round trip failed: input %q, got %q", tt.input, decoded)
			}
		})
	}
}
