package hashutil

import (
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		wantErr  bool
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantErr:  false,
		},
		{
			name:     "simple string",
			input:    []byte("hello"),
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			wantErr:  false,
		},
		{
			name:     "hello world",
			input:    []byte("hello world"),
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
			wantErr:  false,
		},
		{
			name:     "special characters",
			input:    []byte("!@#$%^&*()"),
			expected: "95ce789c5c9d18490972709838ca3a9719094bca3ac16332cfec0652b0236141",
			wantErr:  false,
		},
		{
			name:     "unicode characters",
			input:    []byte("hello 世界"),
			expected: "2e2625f7c51b4a2c75274ab307e86411f57aab475f4a4078df53533f7771bc7f",
			wantErr:  false,
		},
		{
			name:     "newline characters",
			input:    []byte("line1\nline2\nline3"),
			expected: "6bb6a5ad9b9c43a7cb535e636578716b64ac42edea814a4cad102ba404946837",
			wantErr:  false,
		},
		{
			name:     "large input",
			input:    []byte(strings.Repeat("a", 10000)),
			expected: "27dd1f61b867b6a0f6e9d8a41c43231de52107e53ae424de8f847b821db4b711",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("Hash() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHashConsistency(t *testing.T) {
	input := []byte("consistency test")

	hash1, err := Hash(input)
	if err != nil {
		t.Fatalf("first Hash() call failed: %v", err)
	}

	hash2, err := Hash(input)
	if err != nil {
		t.Fatalf("second Hash() call failed: %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("Hash() is not consistent: first=%s, second=%s", hash1, hash2)
	}
}

func BenchmarkHash(b *testing.B) {
	input := []byte("benchmark test input")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Hash(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHashLarge(b *testing.B) {
	input := []byte(strings.Repeat("a", 1000000))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Hash(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
