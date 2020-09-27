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

	return T{float32(2. / rml), 0, 0, 0, 0, float32(2. / tmb), 0, 0, 0, 0, float32(-2. / fmn), 0, float32(-(right + left) / rml), float32(-(top + bottom) / tmb), float32(-(far + near) / fmn), 1}
}

// Perspective generates a perspective projection matrix.
func Perspective(fovy, aspect, near, far float32) T {
	// fovy = (fovy * math.Pi) / 180.0 // convert from degrees to radians
	nmf, f := near-far, float32(1./math.Tan(float64(fovy)/2.0))

	return T{float32(f / aspect), 0, 0, 0, 0, float32(f), 0, 0, 0, 0, float32((near + far) / nmf), -1, 0, 0, float32((2. * far * near) / nmf), 0}
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
