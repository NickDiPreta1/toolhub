package fileconvert

import (
	"bytes"
	"io"
)

func ToUpperText(r io.Reader) (io.Reader, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	upper := bytes.ToUpper(data)

	return bytes.NewReader(upper), nil
}
