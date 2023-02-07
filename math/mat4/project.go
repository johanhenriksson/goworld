package mat4

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Orthographic generates a left-handed orthographic projection matrix.
// Outputs depth values in the range [0, 1]
func Orthographic(left, right, bottom, top, near, far float32) T {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)
	return T{
		2 / rml, 0, 0, 0,
		0, 2 / tmb, 0, 0,
		0, 0, 1 / fmn, 0,
		(right + left) / rml,
		-(top + bottom) / tmb,
		-near / fmn,
		1,
	}
}

// OrthographicRZ generates a left-handed orthographic projection matrix.
// Outputs depth values in the range [1, 0] (reverse Z)
func OrthographicRZ(left, right, bottom, top, near, far float32) T {
	rml, tmb, fmn := (right - left), (top - bottom), (near - far)

	return T{
		2 / rml, 0, 0, 0,
		0, 2 / tmb, 0, 0,
		0, 0, 1 / fmn, 0,
		-(right + left) / rml,
		-(top + bottom) / tmb,
		near / fmn,
		1,
	}
}

// Perspective generates a left-handed perspective projection matrix with reversed depth.
// Outputs depth values in the range [0, 1]
func Perspective(fovy, aspect, near, far float32) T {
	fovy = math.DegToRad(fovy)
	tanHalfFov := math.Tan(fovy) / 2

	return T{
		1 / (aspect * tanHalfFov), 0, 0, 0,
		0, -1 / tanHalfFov, 0, 0,
		0, 0, far / (far - near), 1,
		0, 0, -(far * near) / (far - near), 0,
	}
}

func LookAt(eye, center, up vec3.T) T {
	f := center.Sub(eye).Normalized()
	r := vec3.Cross(up, f).Normalized()
	u := vec3.Cross(f, r)

	M := T{
		r.X, u.X, f.X, 0,
		r.Y, u.Y, f.Y, 0,
		r.Z, u.Z, f.Z, 0,
		0, 0, 0, 1,
	}

	et := Translate(eye.Scaled(-1))
	return M.Mul(&et)
}
