package render

/* Interface for data types that can be uploaded into a vertex buffer object */
type VertexData interface {
    Elements()  int     /* Number of items, usually len(slice) */
    Size()      int     /* Size of each individual element */
}

/* Simple implementation of the VertexData interface for buffering arrays of 32-bit floats */
type FloatBuffer []float32

func (vtx FloatBuffer) Elements() int {
    return len(vtx)
}

func (vtx FloatBuffer) Size() int {
    return 4
}
