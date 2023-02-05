package device

import "unsafe"

func Memcpy(dst, src unsafe.Pointer, n int) int {
	copy(unsafe.Slice((*byte)(dst), n), unsafe.Slice((*byte)(src), n))
	return n
}
