package byte4

// T is a 4-component vector of uint8 (bytes)
type T struct {
	X, Y, Z, W byte
}

func New(x, y, z, w byte) T {
	return T{x, y, z, w}
}
