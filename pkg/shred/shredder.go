package shred

import (
	"bytes"
	"fmt"
)

func Concat(shreds []Shred) ([]byte, error) {
	var buf bytes.Buffer
	for _, shred := range shreds {
		data, ok := shred.Data()
		if !ok {
			return nil, fmt.Errorf("invalid data shred")
		}
		buf.Write(data)
	}
	return buf.Bytes(), nil
}
