package geometry

// a "default" vertex for the "default" shader

type DefaultVertex struct {
	X, Y, Z    float32
	U, V       float32
	Nx, Ny, Nz float32
}

type DefaultVertices []DefaultVertex

func (buffer DefaultVertices) Elements() int {
	return len(buffer)
}

func (buffer DefaultVertices) Size() int {
	return 32
}
