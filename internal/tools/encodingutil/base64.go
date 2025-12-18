package encodingutil

import "encoding/base64"

func Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Decode(input string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}

	return string(decodedBytes), nil
}
