package render

// ScreenQuad is a clip space quad for full screen renders (F32_XYZUV)
var ScreenQuad = FloatBuffer{
	-1, -1, 0, 0, 0,
	1, 1, 0, 1, 1,
	-1, 1, 0, 0, 1,

	-1, -1, 0, 0, 0,
	1, -1, 0, 1, 0,
	1, 1, 0, 1, 1,
}

// Quad is a drawable quad
type Quad struct {
	vao *VertexArray
	mat *Material
}

// NewQuad creates a new quad with a given material
func NewQuad(mat *Material) *Quad {
	q := &Quad{
		vao: CreateVertexArray(Triangles, "geometry"),
		mat: mat,
	}
	q.vao.Buffer("geometry", ScreenQuad)
	if mat != nil {
		mat.SetupVertexPointers()
	}
	return q
}

// Draw the quad
func (q *Quad) Draw() {
	q.vao.Bind()
	if q.mat != nil {
		q.mat.Use()
	}

	q.vao.Draw()
}
