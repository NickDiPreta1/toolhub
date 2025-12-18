package encodingutil

import "encoding/base64"

// Encode returns the base64 encoding of the input string.
func Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// Decode converts a base64 string back to plain text.
func Decode(input string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}

	return string(decodedBytes), nil
}
