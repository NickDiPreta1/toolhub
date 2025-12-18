package fileconvert

import (
	"bytes"
	"io"
)

// ToUpperText reads all input and returns a reader for the uppercased content.
func ToUpperText(r io.Reader) (io.Reader, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	upper := bytes.ToUpper(data)

	return bytes.NewReader(upper), nil
}
