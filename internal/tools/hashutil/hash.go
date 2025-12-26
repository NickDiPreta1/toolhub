package hashutil

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(input []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(input)
	if err != nil {
		return "", err
	}

	hashedBytes := h.Sum(nil)
	hexed := hex.EncodeToString(hashedBytes)

	return hexed, nil
}
