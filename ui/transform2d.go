package ui

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"math"
)

type Transform2D struct {
	Matrix   mgl.Mat4
	Position mgl.Vec3
	Scale    mgl.Vec2
	Rotation float32
}

/* Creates a new 2D transform */
func CreateTransform2D(x, y, z float32) *Transform2D {
	t := &Transform2D{
		Matrix:   mgl.Ident4(),
		Position: mgl.Vec3{x, y, z},
		Scale:    mgl.Vec2{1, 1},
		Rotation: 0.0,
	}
	t.Update(0)
	return t
}

func (t *Transform2D) Update(dt float32) {
	/* Update transform */
	rad := t.Rotation * math.Pi / 180.0
	rotation := mgl.AnglesToQuat(0, 0, rad, mgl.XYZ).Mat4()
	scaling := mgl.Scale3D(t.Scale[0], t.Scale[1], 1)
	translation := mgl.Translate3D(t.Position[0], t.Position[1], t.Position[2])

	/* New transform matrix: S * R * T */
	t.Matrix = scaling.Mul4(rotation.Mul4(translation))
}
