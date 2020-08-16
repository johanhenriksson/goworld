package engine

import "unsafe"

// MeshData is freeform float vertex data
type MeshData struct {
	Items  int
	Buffer []float32
}

// Size returns the byte size of a single element
func (v *MeshData) Size() int {
	return 4 * len(v.Buffer) / v.Items
}

// Elements returns the number of elements
func (v *MeshData) Elements() int {
	return v.Items
}

// Pointer returns an unsafe pointer to the buffer
func (v *MeshData) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&v.Buffer[0])
}
