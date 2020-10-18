package vec4

import "unsafe"

// Array holds an array of 4-component vectors
type Array []T

// Elements returns the number of elements in the array
func (a Array) Elements() int {
	return len(a)
}

// Size return the byte size of an element
func (a Array) Size() int {
	return 16
}

// Pointer returns an unsafe pointer to the first element in the array
func (a Array) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&a[0])
}
