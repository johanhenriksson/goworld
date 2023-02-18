package vertex

type Args interface{}

type GeneratedMesh[A Args, V Vertex, I Index] interface {
	Mesh
	Update(A)
}

type generated[A Args, V Vertex, I Index] struct {
	Mesh
	key       string
	version   int
	hash      int
	generator func(A) (V, I)
}

func NewGenerated[A Args, V Vertex, I Index](key string, args A, generator func(A) (V, I)) GeneratedMesh[A, V, I] {
	return &generated[A, V, I]{
		key:       key,
		version:   1,
		generator: generator,
	}
}

func (g *generated[A, V, I]) Update(args A) {
	// if args hash has changed, update version
}
