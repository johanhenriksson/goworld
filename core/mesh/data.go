package mesh

import (
	"unsafe"
)

type Data interface {
	Size() int
	Elements() int
	Pointer() unsafe.Pointer
}

func NewData(items int, buffer []float32) Data {
	return &data{
		items:  items,
		buffer: buffer,
	}
}

// data is freeform float vertex data
type data struct {
	items  int
	buffer []float32
}

// Size returns the byte size of a single element
func (v *data) Size() int {
	return 4 * len(v.buffer) / v.items
}

// Elements returns the number of elements
func (v *data) Elements() int {
	return v.items
}

// Pointer returns an unsafe pointer to the buffer
func (v *data) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&v.buffer[0])
}
