package render

/* Clip space quad vertex data
 * float32 X Y Z U V */
var ScreenQuad FloatBuffer = FloatBuffer{
	-1, -1, 0, 0, 0,
	1, 1, 0, 1, 1,
	-1, 1, 0, 0, 1,

	-1, -1, 0, 0, 0,
	1, -1, 0, 1, 0,
	1, 1, 0, 1, 1,
}

type RenderQuad struct {
	vao *VertexArray
	vbo *VertexBuffer
}

func NewRenderQuad() *RenderQuad {
	q := &RenderQuad{
		vao: CreateVertexArray(),
		vbo: CreateVertexBuffer(),
	}
	q.vao.Length = 6 // two triangles, six vertices
	q.vao.Bind()
	q.vbo.Buffer(ScreenQuad)
	return q
}

func (q *RenderQuad) Draw() {
	q.vao.Draw()
}
