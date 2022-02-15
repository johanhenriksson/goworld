package mat4

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"

	mgl "github.com/go-gl/mathgl/mgl32"
)

// Orthographic generates a right-handed orthographic projection matrix.
// Outputs depth values in the range [-1, 1]
func Orthographic(left, right, bottom, top, near, far float32) T {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)

	return T{
		2 / rml, 0, 0, 0,
		0, 2 / tmb, 0, 0,
		0, 0, -2 / fmn, 0,
		-(right + left) / rml, -(top + bottom) / tmb, -(far + near) / fmn, 1,
	}
}

// OrthographicLH generates a left-handed orthographic projection matrix.
// Outputs depth values in the range [-1, 1]
func OrthographicLH(left, right, bottom, top, near, far float32) T {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)

	return T{
		2 / rml, 0, 0, 0,
		0, 2 / tmb, 0, 0,
		0, 0, 2 / fmn, 0,
		-(right + left) / rml, -(top + bottom) / tmb, -(far + near) / fmn, 1,
	}
}

// OrthographicVK generates a left-handed orthographic projection matrix.
// Outputs depth values in the range [1, 0]
func OrthographicVK(left, right, bottom, top, near, far float32) T {
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

// Perspective generates a right-handed perspective projection matrix.
// Outputs depth in the range [-1, 1]
func Perspective(fovy, aspect, near, far float32) T {
	fovy = math.DegToRad(fovy)
	fmn, tanHalfFov := far-near, math.Tan(fovy/2)

	return T{
		1 / (aspect * tanHalfFov), 0, 0, 0,
		0, 1 / tanHalfFov, 0, 0,
		0, 0, -(near + far) / fmn, -1,
		0, 0, (2 * far * near) / fmn, 0,
	}
}

// PerspectiveLH generates a left-handed perspective projection matrix.
// Outputs depth in the range [-1, 1]
func PerspectiveLH(fovy, aspect, near, far float32) T {
	fovy = math.DegToRad(fovy)
	tanHalfFov := math.Tan(fovy) / 2

	return T{
		1 / (aspect * tanHalfFov), 0, 0, 0,
		0, 1 / tanHalfFov, 0, 0,
		0, 0, -(far + near) / (far - near), 1,
		0, 0, (2 * far * near) / (far - near), 0,
	}
}

// PerspectiveVK generates a left-handed perspective projection matrix with reversed depth.
// Outputs depth in the range [1, 0]
func PerspectiveVK(fovy, aspect, near, far float32) T {
	fovy = math.DegToRad(fovy)
	tanHalfFov := math.Tan(fovy) / 2

	return T{
		1 / (aspect * tanHalfFov), 0, 0, 0,
		0, -1 / tanHalfFov, 0, 0,
		0, 0, far / (far - near), 1,
		0, 0, -(far * near) / (far - near), 0,
	}
}

// LookAt generates a transform matrix from world space into the specific eye
// space.
func LookAt(eye, center vec3.T) T {
	mat := mgl.LookAtV(
		mgl.Vec3{eye.X, eye.Y, eye.Z},
		mgl.Vec3{center.X, center.Y, center.Z},
		mgl.Vec3{0, 1, 0},
	)
	return T(mat)
}

func LookAtLH(eye, center vec3.T) T {
	up := vec3.UnitY
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
