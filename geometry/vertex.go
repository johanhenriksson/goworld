package geometry

import "unsafe"

// Vertex is a typical vertex with positions, uvs, and normals
type Vertex struct {
	X, Y, Z    float32
	U, V       float32
	Nx, Ny, Nz float32
}

// DefaultVertices is a GPU bufferable slice of default vertices
type Vertices []Vertex

func (buffer Vertices) Elements() int {
	return len(buffer)
}

func (buffer Vertices) Size() int {
	return 32
}

func (buffer Vertices) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&buffer[0])
}
