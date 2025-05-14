package server

import (
	"crypto/rand"
	"encoding/binary"
)

func generateURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 8
	b := make([]byte, length)

	for i := range b {
		var n uint64
		binary.Read(rand.Reader, binary.LittleEndian, &n)
		b[i] = charset[n%uint64(len(charset))]
	}
	return "room-" + string(b)
}
