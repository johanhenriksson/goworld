package util

import (
    "C"
    "unsafe"
    "strings"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func GLString(str string) *uint8 {
    if !strings.HasSuffix(str, "\x00") {
        str += "\x00"
    }
    return gl.Str(str)
}

func GoString(cstring *uint8) string {
    return C.GoString((*C.char)(unsafe.Pointer(cstring)))
}
