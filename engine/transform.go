package engine

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"math"
)

/* Represents a 3D transformation */
type Transform struct {
	Matrix   mgl.Mat4
	Position mgl.Vec3
	Rotation mgl.Vec3
	Scale    mgl.Vec3
	Forward  mgl.Vec3
	Right    mgl.Vec3
	Up       mgl.Vec3
	/* Probably needs a changed flag */
}

/* Creates a new 3D transform */
func CreateTransform(x, y, z float32) *Transform {
	t := &Transform{
		Matrix:   mgl.Ident4(),
		Position: mgl.Vec3{x, y, z},
		Rotation: mgl.Vec3{0, 0, 0},
		Scale:    mgl.Vec3{1, 1, 1},
	}
	t.Update(0)
	return t
}

/* Update transform matrix and right/up/forward vectors */
func (t *Transform) Update(dt float32) {
	// todo: avoid recalculating unless something has changed

	/* Update transform */
	rad := t.Rotation.Mul(math.Pi / 180.0) // translate rotaiton to radians
	rotation := mgl.AnglesToQuat(rad[0], rad[1], rad[2], mgl.XYZ).Mat4()
	scaling := mgl.Scale3D(t.Scale[0], t.Scale[1], t.Scale[2])
	translation := mgl.Translate3D(t.Position[0], t.Position[1], t.Position[2])

	/* New transform matrix: S * R * T */
	//m := scaling.Mul4(rotation.Mul4(translation))
	m := translation.Mul4(rotation.Mul4(scaling))

	/* Grab axis vectors from transformation matrix */
	t.Right[0] = m[4*0+0] // first column
	t.Right[1] = m[4*1+0]
	t.Right[2] = m[4*2+0]
	t.Up[0] = m[4*0+1] // second column
	t.Up[1] = m[4*1+1]
	t.Up[2] = m[4*2+1]
	t.Forward[0] = -m[4*0+2] // third column
	t.Forward[1] = -m[4*1+2]
	t.Forward[2] = -m[4*2+2]

	/* Update transformation matrix */
	t.Matrix = m
}

func (t *Transform) Translate(offset mgl.Vec3) {
	t.Position = t.Position.Add(offset)
}

/* Transforms a point into this coordinate system */
func (t *Transform) TransformPoint(point mgl.Vec3) mgl.Vec3 {
	p4 := mgl.Vec4{point[0], point[1], point[2], 1}
	return t.Matrix.Mul4x1(p4).Vec3()
}

/* Transforms a direction into this coordinate system */
func (t *Transform) TransformDir(dir mgl.Vec3) mgl.Vec3 {
	d4 := mgl.Vec4{dir[0], dir[1], dir[2], 0}
	return t.Matrix.Mul4x1(d4).Vec3()
}

// todo: InverseTransformPoint
// todo: InverseTransformDir
