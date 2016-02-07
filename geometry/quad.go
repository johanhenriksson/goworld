package geometry

import (
    "math"
    "github.com/johanhenriksson/goworld/render"
)

type Quad struct {
    TopLeft     ColorVertex
    TopRight    ColorVertex
    BottomLeft  ColorVertex
    BottomRight ColorVertex
    Material    *render.Material

    segments    int
    border      float64
    vao         *render.VertexArray
    vbo         *render.VertexBuffer
}

func NewQuad(mat *render.Material, color render.Color, w,h,z float32) *Quad {
    q := &Quad {
        Material:    mat,
        TopLeft:     ColorVertex { X: 0, Y: h, Z: z, Color: color, },
        TopRight:    ColorVertex { X: w, Y: h, Z: z, Color: color, },
        BottomLeft:  ColorVertex { X: 0, Y: 0, Z: z, Color: color, },
        BottomRight: ColorVertex { X: w, Y: 0, Z: z, Color: color, },
        segments:    5,
        border:      0,

        vao: render.CreateVertexArray(),
        vbo: render.CreateVertexBuffer(),
    }
    q.compute()
    return q
}

func (q *Quad) BorderWidth() float64 { return q.border }
func (q *Quad) SetBorderWidth(width float64) {
    q.border = width
    q.compute()
}

func (q *Quad) appendCorner(vtx *ColorVertices, origin ColorVertex, n int, r, offset float64) {
    if n == 0 {
        /* TODO: Square corner */
    } else {
        /* Rounded corner */
        v := (math.Pi / 2.0) / float64(n)
        var prev ColorVertex
        for i := 0; i <= n; i++ {
            p := ColorVertex {
                X: origin.X + float32(r * math.Cos(offset + float64(i)*v)),
                Y: origin.Y + float32(r * math.Sin(offset + float64(i)*v)),
                Z: origin.Z,
                Color: origin.Color,
            }

            if i > 0 {
                *vtx = append(*vtx, origin, prev, p)
            }

            prev = p
        }
    }
}

func (q *Quad) compute() {
    vtx := ColorVertices {
        q.BottomLeft, q.TopRight, q.TopLeft,
        q.BottomLeft, q.BottomRight, q.TopRight,
    }

    /* If we have a positive border width, tesselate border */
    b := float32(q.border)
    if b > 0.0 {
        q.appendCorner(&vtx, q.TopRight, q.segments, q.border, 0.0)
        q.appendCorner(&vtx, q.TopLeft, q.segments, q.border, math.Pi/2.0)
        q.appendCorner(&vtx, q.BottomLeft, q.segments, q.border, math.Pi)
        q.appendCorner(&vtx, q.BottomRight, q.segments, q.border, 3.0*math.Pi/2.0)

        /* Top Border Box */
        topTopLeft := q.TopLeft
        topTopLeft.Y += b
        topTopRight := q.TopRight
        topTopRight.Y += b
        vtx = append(vtx, q.TopLeft, topTopRight, topTopLeft,
                          q.TopLeft, q.TopRight, topTopRight)

        /* Bottom border box */
        bottomBottomLeft := q.BottomLeft
        bottomBottomLeft.Y -= b
        bottomBottomRight := q.BottomRight
        bottomBottomRight.Y -= b
        vtx = append(vtx, bottomBottomLeft, q.BottomRight, q.BottomLeft,
                          bottomBottomLeft, bottomBottomRight, q.BottomRight)

        /* Right border box */
        rightTopRight := q.TopRight
        rightTopRight.X += b
        rightBottomRight := q.BottomRight
        rightBottomRight.X += b
        vtx = append(vtx, q.BottomRight, rightTopRight, q.TopRight,
                          q.BottomRight, rightBottomRight, rightTopRight)

        /* Left border box */
        leftTopLeft := q.TopLeft
        leftTopLeft.X -= b
        leftBottomLeft := q.BottomLeft
        leftBottomLeft.X -= b
        vtx = append(vtx, leftBottomLeft, q.TopLeft, leftTopLeft,
                          leftBottomLeft, q.BottomLeft, q.TopLeft)
    }

    /* Setup VAO */
    q.vao.Length = int32(len(vtx))
    q.vao.Bind()
    q.vbo.Buffer(vtx)
    q.Material.SetupVertexPointers()
}

func (q *Quad) Draw(args render.DrawArgs) {
    q.Material.Use()
    q.Material.Shader.Matrix4f("model", &args.Transform[0])
    q.Material.Shader.Matrix4f("viewport", &args.Projection[0])
    q.vao.Draw()
}
