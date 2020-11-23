package mat4

// Ident returns a new 4x4 identity matrix
func Ident() T {
	return T{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}
