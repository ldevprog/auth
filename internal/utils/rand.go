package utils

import (
	"crypto/rand"
	binary "encoding/binary"
)

func RandInt64Positive() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0
	}

	u := int64(binary.LittleEndian.Uint64(b[:]))
	return int64(u & 0x7FFFFFFFFFFFFFFF)
}
