package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// PrettyPrint validates JSON then returns a formatted version.
func PrettyPrint(input string) (string, error) {
	if err := validateJSON(input); err != nil {
		return "", err
	}

	var store interface{}
	if err := json.Unmarshal([]byte(input), &store); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	mBytes, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}

	return string(mBytes), nil
}

// Minify validates JSON then removes whitespace.
func Minify(input string) (string, error) {
	if err := validateJSON(input); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(input)); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	return buf.String(), nil
}

// validateJSON is a small guard for empty input.
func validateJSON(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input cannot be empty")
	}

	return nil
}
