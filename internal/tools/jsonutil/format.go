package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

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

func Minify(input string) (string, error) {
	if err := validateJSON(input); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	json.Compact(&buf, []byte(input))

	return buf.String(), nil
}

func validateJSON(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input cannot be empty")
	}

	return nil
}
