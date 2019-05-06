package geometry

/** ImageVertex */
type ImageVertex struct {
	X, Y, Z float32 // 12 bytes
	Tx, Ty  float32 // 8 bytes
}

type ImageVertices []ImageVertex

func (buffer ImageVertices) Elements() int { return len(buffer) }
func (buffer ImageVertices) Size() int     { return 20 }
