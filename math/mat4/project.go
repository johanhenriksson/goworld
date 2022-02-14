// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat4

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/johanhenriksson/goworld/math/vec3"
)

// Orthographic generates an orthographic projection matrix.
func Orthographic(left, right, bottom, top, near, far float32) T {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)

	return T{
		2 / rml, 0, 0, 0,
		0, 2 / tmb, 0, 0,
		0, 0, -2 / fmn, 0,
		-(right + left) / rml, -(top + bottom) / tmb, -(far + near) / fmn, 1,
	}
}

func OrthographicLH(left, right, bottom, top, near, far float32) T {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)

	return T{
		2 / rml, 0, 0, 0,
		0, 2 / tmb, 0, 0,
		0, 0, 2 / fmn, 0,
		-(right + left) / rml, -(top + bottom) / tmb, -(far + near) / fmn, 1,
	}
}

// #ifdef GLM_DEPTH_ZERO_TO_ONE
// 		Result[2][2] = farVal / (farVal - nearVal);
// 		Result[3][2] = -(farVal * nearVal) / (farVal - nearVal);
// #else
// 		Result[2][2] = (farVal + nearVal) / (farVal - nearVal);
// 		Result[3][2] = - (static_cast<T>(2) * farVal * nearVal) / (farVal - nearVal);
// #endif

// Perspective generates a perspective projection matrix.
func Perspective(fovy, aspect, near, far float32) T {
	fovy = (fovy * math.Pi) / 180.0 // convert from degrees to radians
	nmf, f := near-far, float32(1./math.Tan(float64(fovy)/2))

	return T{float32(f / aspect), 0, 0, 0, 0, float32(f), 0, 0, 0, 0, float32((near + far) / nmf), -1, 0, 0, float32((2. * far * near) / nmf), 0}
}

func PerspectiveLH(fovy, aspect, near, far float32) T {
	fovy = (fovy * math.Pi) / 180.0 // convert from degrees to radians
	tanHalfFov := float32(math.Tan(float64(fovy) / 2))

	return T{
		1 / (aspect * tanHalfFov), 0, 0, 0,
		0, 1 / tanHalfFov, 0, 0,
		0, 0, -(far + near) / (far - near), 1,
		0, 0, (2 * far * near) / (far - near), 0,
	}
}

func PerspectiveVK(fovy, aspect, near, far float32) T {
	fovy = (fovy * math.Pi) / 180.0 // convert from degrees to radians
	tanHalfFov := float32(math.Tan(float64(fovy) / 2))

	return T{
		1 / (aspect * tanHalfFov), 0, 0, 0,
		0, 1 / tanHalfFov, 0, 0,
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
