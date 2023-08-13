package vertex

import "github.com/johanhenriksson/goworld/math/vec3"

type Triangle struct {
	A, B, C vec3.T
}

func (t *Triangle) Normal() vec3.T {
	// Set Vector U to (Triangle.p2 minus Triangle.p1)
	u := t.B.Sub(t.A)
	// Set Vector V to (Triangle.p3 minus Triangle.p1)
	v := t.C.Sub(t.A)

	// Set Normal.x to (multiply U.y by V.z) minus (multiply U.z by V.y)
	x := u.Y*v.Z - u.Z*v.Y
	// Set Normal.y to (multiply U.z by V.x) minus (multiply U.x by V.z)
	y := u.Z*v.X - u.X*v.Z
	// Set Normal.z to (multiply U.x by V.y) minus (multiply U.y by V.x)
	z := u.X*v.Y - u.Y*v.X

	return vec3.New(x, y, z).Normalized()
}
