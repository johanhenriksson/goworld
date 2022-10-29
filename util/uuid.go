package util

import (
	"math/rand"
)

var idCharset = []byte("abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ0123456789")

func NewUUID(length int) string {
	id := make([]byte, length)
	charsetLen := int64(len(idCharset))
	for i := 0; i < length; i++ {
		ch := rand.Int63n(charsetLen)
		id[i] = idCharset[ch]
	}
	return string(id)
}
