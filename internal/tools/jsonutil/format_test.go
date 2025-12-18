package jsonutil

import (
	"encoding/json"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		// Basic valid JSON - keys alphabetically sorted
		{
			name:        "simple object",
			input:       `{"name":"John","age":30}`,
			expected:    "{\n  \"age\": 30,\n  \"name\": \"John\"\n}",
			expectError: false,
		},
		{
			name:        "already pretty",
			input:       "{\n  \"name\": \"John\"\n}",
			expected:    "{\n  \"name\": \"John\"\n}",
			expectError: false,
		},
		{
			name:        "simple array",
			input:       `[1,2,3,4,5]`,
			expected:    "[\n  1,\n  2,\n  3,\n  4,\n  5\n]",
			expectError: false,
		},
		{
			name:        "nested object",
			input:       `{"person":{"name":"John","age":30}}`,
			expected:    "{\n  \"person\": {\n    \"age\": 30,\n    \"name\": \"John\"\n  }\n}",
			expectError: false,
		},
		{
			name:        "nested arrays",
			input:       `{"matrix":[[1,2],[3,4]]}`,
			expected:    "{\n  \"matrix\": [\n    [\n      1,\n      2\n    ],\n    [\n      3,\n      4\n    ]\n  ]\n}",
			expectError: false,
		},

		// Different data types
		{
			name:        "string value",
			input:       `{"message":"hello world"}`,
			expected:    "{\n  \"message\": \"hello world\"\n}",
			expectError: false,
		},
		{
			name:        "number value",
			input:       `{"count":42}`,
			expected:    "{\n  \"count\": 42\n}",
			expectError: false,
		},
		{
			name:        "boolean values",
			input:       `{"active":true,"deleted":false}`,
			expected:    "{\n  \"active\": true,\n  \"deleted\": false\n}",
			expectError: false,
		},
		{
			name:        "null value",
			input:       `{"data":null}`,
			expected:    "{\n  \"data\": null\n}",
			expectError: false,
		},
		{
			name:        "mixed types",
			input:       `{"name":"John","age":30,"active":true,"scores":[95,87,92],"address":null}`,
			expected:    "{\n  \"active\": true,\n  \"address\": null,\n  \"age\": 30,\n  \"name\": \"John\",\n  \"scores\": [\n    95,\n    87,\n    92\n  ]\n}",
			expectError: false,
		},

		// Edge cases - valid JSON
		{
			name:        "empty object",
			input:       `{}`,
			expected:    "{}",
			expectError: false,
		},
		{
			name:        "empty array",
			input:       `[]`,
			expected:    "[]",
			expectError: false,
		},
		{
			name:        "single string",
			input:       `"hello"`,
			expected:    "\"hello\"",
			expectError: false,
		},
		{
			name:        "single number",
			input:       `42`,
			expected:    "42",
			expectError: false,
		},
		{
			name:        "single boolean",
			input:       `true`,
			expected:    "true",
			expectError: false,
		},
		{
			name:        "single null",
			input:       `null`,
			expected:    "null",
			expectError: false,
		},

		// Special characters in strings
		{
			name:        "string with quotes",
			input:       `{"message":"He said \"hello\""}`,
			expected:    "{\n  \"message\": \"He said \\\"hello\\\"\"\n}",
			expectError: false,
		},
		{
			name:        "string with newlines",
			input:       `{"message":"line1\nline2"}`,
			expected:    "{\n  \"message\": \"line1\\nline2\"\n}",
			expectError: false,
		},
		{
			name:        "unicode characters",
			input:       `{"emoji":"👋","accent":"café"}`,
			expected:    "{\n  \"accent\": \"café\",\n  \"emoji\": \"👋\"\n}",
			expectError: false,
		},

		// Complex real-world examples - keys alphabetically sorted
		{
			name:        "API response structure",
			input:       `{"status":"success","data":{"id":123,"name":"Product","price":29.99},"meta":{"timestamp":"2024-01-01T00:00:00Z"}}`,
			expected:    "{\n  \"data\": {\n    \"id\": 123,\n    \"name\": \"Product\",\n    \"price\": 29.99\n  },\n  \"meta\": {\n    \"timestamp\": \"2024-01-01T00:00:00Z\"\n  },\n  \"status\": \"success\"\n}",
			expectError: false,
		},

		// Error cases - invalid JSON
		{
			name:        "empty input",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			input:       "   ",
			expected:    "",
			expectError: true,
		},
		{
			name:        "missing closing brace",
			input:       `{"name":"John"`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "missing closing bracket",
			input:       `[1,2,3`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "trailing comma",
			input:       `{"name":"John",}`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "unquoted key",
			input:       `{name:"John"}`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "single quotes instead of double",
			input:       `{'name':'John'}`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "plain text",
			input:       `this is not JSON`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "incomplete object",
			input:       `{"name":}`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "extra closing brace",
			input:       `{"name":"John"}}`,
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PrettyPrint(tt.input)

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("PrettyPrint(%q) expected error but got nil", tt.input)
				}
				return
			}

			// No error expected
			if err != nil {
				t.Errorf("PrettyPrint(%q) unexpected error: %v", tt.input, err)
				return
			}

			if result != tt.expected {
				t.Errorf("PrettyPrint(%q)\ngot:\n%s\n\nwant:\n%s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMinify(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		// Basic valid JSON - keys preserve original order (Compact doesn't re-marshal)
		{
			name:        "pretty object to minified",
			input:       "{\n  \"name\": \"John\",\n  \"age\": 30\n}",
			expected:    `{"name":"John","age":30}`,
			expectError: false,
		},
		{
			name:        "already minified",
			input:       `{"name":"John","age":30}`,
			expected:    `{"name":"John","age":30}`,
			expectError: false,
		},
		{
			name:        "pretty array to minified",
			input:       "[\n  1,\n  2,\n  3\n]",
			expected:    `[1,2,3]`,
			expectError: false,
		},
		{
			name:        "nested object",
			input:       "{\n  \"person\": {\n    \"name\": \"John\"\n  }\n}",
			expected:    `{"person":{"name":"John"}}`,
			expectError: false,
		},

		// Various whitespace patterns
		{
			name:        "extra spaces",
			input:       `{  "name"  :  "John"  }`,
			expected:    `{"name":"John"}`,
			expectError: false,
		},
		{
			name:        "tabs and newlines",
			input:       "{\n\t\"name\": \"John\"\n}",
			expected:    `{"name":"John"}`,
			expectError: false,
		},
		{
			name:        "multiple newlines",
			input:       "{\n\n\n  \"name\": \"John\"\n\n\n}",
			expected:    `{"name":"John"}`,
			expectError: false,
		},

		// Different data types
		{
			name:        "string value",
			input:       `{ "message" : "hello world" }`,
			expected:    `{"message":"hello world"}`,
			expectError: false,
		},
		{
			name:        "number value",
			input:       `{ "count" : 42 }`,
			expected:    `{"count":42}`,
			expectError: false,
		},
		{
			name:        "boolean values",
			input:       `{ "active" : true , "deleted" : false }`,
			expected:    `{"active":true,"deleted":false}`,
			expectError: false,
		},
		{
			name:        "null value",
			input:       `{ "data" : null }`,
			expected:    `{"data":null}`,
			expectError: false,
		},

		// Edge cases - valid JSON
		{
			name:        "empty object with whitespace",
			input:       `{   }`,
			expected:    `{}`,
			expectError: false,
		},
		{
			name:        "empty array with whitespace",
			input:       `[   ]`,
			expected:    `[]`,
			expectError: false,
		},
		{
			name:        "single string with whitespace",
			input:       `  "hello"  `,
			expected:    `"hello"`,
			expectError: false,
		},
		{
			name:        "single number with whitespace",
			input:       `  42  `,
			expected:    `42`,
			expectError: false,
		},

		// Strings with spaces (should be preserved)
		{
			name:        "preserve spaces in string values",
			input:       `{ "message" : "hello   world" }`,
			expected:    `{"message":"hello   world"}`,
			expectError: false,
		},
		{
			name:        "preserve newlines in string values",
			input:       `{ "message" : "line1\nline2" }`,
			expected:    `{"message":"line1\nline2"}`,
			expectError: false,
		},

		// Complex structures - keys preserve original order
		{
			name:        "API response",
			input:       "{\n  \"status\": \"success\",\n  \"data\": {\n    \"id\": 123\n  }\n}",
			expected:    `{"status":"success","data":{"id":123}}`,
			expectError: false,
		},
		{
			name:        "array of objects",
			input:       "[\n  { \"id\": 1 },\n  { \"id\": 2 }\n]",
			expected:    `[{"id":1},{"id":2}]`,
			expectError: false,
		},

		// Error cases - invalid JSON
		{
			name:        "empty input",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			input:       "   \n\t  ",
			expected:    "",
			expectError: true,
		},
		{
			name:        "missing closing brace",
			input:       `{"name":"John"`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "trailing comma",
			input:       `{"name":"John",}`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "unquoted key",
			input:       `{name:"John"}`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "plain text",
			input:       `not json at all`,
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Minify(tt.input)

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("Minify(%q) expected error but got nil", tt.input)
				}
				return
			}

			// No error expected
			if err != nil {
				t.Errorf("Minify(%q) unexpected error: %v", tt.input, err)
				return
			}

			if result != tt.expected {
				t.Errorf("Minify(%q)\ngot:  %s\nwant: %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestPrettyPrintMinifyRoundTrip tests that pretty printing then minifying
// produces valid, compact JSON. Note: Minify preserves key order from input.
func TestPrettyPrintMinifyRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string // Start with minified JSON
	}{
		{"simple object", `{"name":"John","age":30}`},
		{"nested object", `{"person":{"name":"John","address":{"city":"NYC"}}}`},
		{"array", `[1,2,3,4,5]`},
		{"array of objects", `[{"id":1},{"id":2}]`},
		{"mixed types", `{"string":"text","number":42,"bool":true,"null":null}`},
		{"empty object", `{}`},
		{"empty array", `[]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Pretty print the minified input
			pretty, err := PrettyPrint(tt.input)
			if err != nil {
				t.Fatalf("PrettyPrint failed: %v", err)
			}

			// Step 2: Minify the pretty-printed result
			minified, err := Minify(pretty)
			if err != nil {
				t.Fatalf("Minify failed: %v", err)
			}

			// Step 3: Parse both to verify they're semantically equal
			// (Key order may differ, but content should be the same)
			var original, result interface{}
			if err := json.Unmarshal([]byte(tt.input), &original); err != nil {
				t.Fatalf("Failed to parse original: %v", err)
			}
			if err := json.Unmarshal([]byte(minified), &result); err != nil {
				t.Fatalf("Failed to parse result: %v", err)
			}

			// Compare by re-marshaling both (this sorts keys consistently)
			originalBytes, _ := json.Marshal(original)
			resultBytes, _ := json.Marshal(result)

			if string(originalBytes) != string(resultBytes) {
				t.Errorf("Round trip produced different JSON:\noriginal: %s\nresult:   %s", originalBytes, resultBytes)
			}
		})
	}
}

// TestMinifyPrettyPrintRoundTrip tests that minifying then pretty printing
// preserves JSON structure (though PrettyPrint will sort keys)
func TestMinifyPrettyPrintRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string // Start with pretty JSON
	}{
		{
			name:  "simple object",
			input: "{\n  \"name\": \"John\"\n}",
		},
		{
			name:  "object with multiple keys",
			input: "{\n  \"name\": \"John\",\n  \"age\": 30\n}",
		},
		{
			name:  "nested object",
			input: "{\n  \"person\": {\n    \"name\": \"John\"\n  }\n}",
		},
		{
			name:  "array",
			input: "[\n  1,\n  2,\n  3\n]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Minify the pretty input
			minified, err := Minify(tt.input)
			if err != nil {
				t.Fatalf("Minify failed: %v", err)
			}

			// Step 2: Pretty print the minified result
			pretty, err := PrettyPrint(minified)
			if err != nil {
				t.Fatalf("PrettyPrint failed: %v", err)
			}

			// Step 3: Verify both are valid and semantically equal
			var original, result interface{}
			if err := json.Unmarshal([]byte(tt.input), &original); err != nil {
				t.Fatalf("Failed to parse original: %v", err)
			}
			if err := json.Unmarshal([]byte(pretty), &result); err != nil {
				t.Fatalf("Failed to parse result: %v", err)
			}

			// Compare by re-marshaling (sorts keys consistently)
			originalBytes, _ := json.Marshal(original)
			resultBytes, _ := json.Marshal(result)

			if string(originalBytes) != string(resultBytes) {
				t.Errorf("Round trip produced different JSON:\noriginal: %s\nresult:   %s", originalBytes, resultBytes)
			}
		})
	}
}
