package render

// VertexData is an interface for data types that can be stored in a vertex buffer object
type VertexData interface {
	Elements() int /* Number of items, usually len(slice) */
	Size() int     /* Size of each individual element */
}

// FloatBuffer is a simple implementation of the VertexData interface for buffering arrays of 32-bit floats
type FloatBuffer []float32

// Elements returns the number of vertex elements in the buffer
func (vtx FloatBuffer) Elements() int {
	return len(vtx)
}

// Size returns the byte size of a buffer element
func (vtx FloatBuffer) Size() int {
	return 4
}
