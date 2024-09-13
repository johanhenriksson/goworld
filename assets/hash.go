package assets

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"unsafe"
)

func Hash(v any) string {
	rv := reflect.ValueOf(v)

	// Handle nil values
	if !rv.IsValid() {
		panic("cannot hash nil value")
	}

	// Dereference pointers and interfaces
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	var ptr unsafe.Pointer
	var size uintptr

	switch rv.Kind() {
	case reflect.UnsafePointer:
		panic("cannot hash unsafe.Pointer")
	case reflect.Slice, reflect.Array:
		ptr = unsafe.Pointer(rv.Pointer())
		if rv.Type().Kind() == reflect.Ptr {
			panic("cannot hash pointer slice")
		}
		size = uintptr(rv.Len()) * rv.Type().Size()
	case reflect.String:
		str := rv.String()
		ptr = unsafe.Pointer(unsafe.StringData(str))
		size = uintptr(len(str))
	default:
		// for other types, get a pointer to the value
		ptr = unsafe.Pointer(rv.UnsafeAddr())
		size = rv.Type().Size()
	}

	// create a byte slice that represents the memory without copying
	buf := unsafe.Slice((*byte)(ptr), size)

	// calculate SHA-256 hash
	hash := sha256.Sum256(buf)
	return hex.EncodeToString(hash[:])
}
