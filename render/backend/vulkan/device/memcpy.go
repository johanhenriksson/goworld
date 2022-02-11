package device

//#include <string.h>
import "C"
import "unsafe"

func memcpy(dst, src unsafe.Pointer, n int) int {
	C.memcpy(dst, src, C.size_t(n))
	return n
}
