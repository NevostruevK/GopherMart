package fgzip

import (
	"bytes"
	"compress/gzip"
)

func Decompress(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(gz)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
