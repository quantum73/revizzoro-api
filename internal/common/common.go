package common

import (
	"bytes"
)

func ConcatenateBytesSlice(data [][]byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteRune('[')
	for _, d := range data {
		buf.Write(d)
		buf.WriteRune(',')
	}
	buf.WriteRune(']')

	return buf.Bytes(), nil
}
