package object

import (
	"math/rand"
	"strconv"
)

func Key(prefix string, object Component) string {
	p := len(prefix)
	buffer := make([]byte, p+1, p+9)
	copy(buffer, []byte(prefix))
	buffer[p] = '-'
	dst := strconv.AppendUint(buffer, uint64(object.ID()), 10)
	return string(dst)
}

func ID() uint {
	return uint(rand.Int63n(0xFFFFFFFF))
}
