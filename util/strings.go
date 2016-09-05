package util

import (
    "C"
    "unsafe"
    "strings"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// Converts a Go string to a uint8 array for use with OpenGL.
// Returns a pointer to the char array and a function to free the memory associated with the array
func GLString(str string) (**uint8, func()) {
    if !strings.HasSuffix(str, "\x00") {
        str += "\x00"
    }
    return gl.Strs(str)
}

// Converts a C string to a Go string
func GoString(cstring *uint8) string {
    return C.GoString((*C.char)(unsafe.Pointer(cstring)))
}
