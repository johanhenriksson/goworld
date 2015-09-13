package ui
import (
    "fmt"
    "math"
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/geometry"
)

type Quad struct {
    TopLeft     ColorVertex
    TopRight    ColorVertex
    BottomLeft  ColorVertex
    BottomRight ColorVertex
    Material    *render.Material
    vao         *geometry.VertexArray
    vbo         *geometry.VertexBuffer
}

func NewQuad(mat *render.Material, color Color, w,h,z,r,g,b,a float32) *Quad {
    q := &Quad {
        Material: mat,
        TopLeft: ColorVertex {
            X: 0,
            Y: h,
            Z: z,
            Color: color,
        },
        TopRight: ColorVertex {
            X: w,
            Y: h,
            Z: z,
            Color: color,
        },
        BottomLeft: ColorVertex {
            X: 0,
            Y: 0,
            Z: z,
            Color: color,
        },
        BottomRight: ColorVertex {
            X: w,
            Y: 0,
            Z: z,
            Color: color,
        },
        vao: geometry.CreateVertexArray(),
        vbo: geometry.CreateVertexBuffer(),
    }
    q.compute()
    return q
}

func (q *Quad) compute() {
    vtx := ColorVertices {
        q.BottomLeft, q.TopRight, q.TopLeft,
        q.BottomLeft, q.BottomRight, q.TopRight,
    }

    n := 5
    v := (math.Pi / 2.0) / float64(n)
    b := float32(50.0)
    origin := q.TopRight
    prev := ColorVertex {
        X: origin.X, Y: origin.Y + b, Z: origin.Z,
        Color: origin.Color,
    }
    fmt.Println("Vtx:", len(vtx))
    for i := 1; i <= n; i++ {
        p := ColorVertex {
            X: origin.X + b * float32(math.Cos((math.Pi/2.0) - float64(i)*v)),
            Y: origin.Y + b * float32(math.Sin((math.Pi/2.0) - float64(i)*v)),
            Z: origin.Z,
            Color: origin.Color,
        }
        fmt.Println("---- Triangle", i, "----")
        fmt.Println(origin)
        fmt.Println(p)
        fmt.Println(prev)
        vtx = append(vtx, origin, p, prev)
        prev = p
    }
    fmt.Println("Vtx:", len(vtx))

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

    q.vao.Length = int32(len(vtx))
    q.vao.Bind()
    q.vbo.Buffer(vtx)
    q.Material.Setup()
}

func (q *Quad) Draw(args DrawArgs) {
    q.Material.Use()
    q.Material.Shader.Matrix4f("model", &args.Transform[0])
    q.Material.Shader.Matrix4f("viewport", &args.Viewport[0])
    q.vao.Draw()
}
