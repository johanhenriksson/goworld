// Based on code from github.com/go-gl/mathgl:
// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package mat4

import (
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Translate returns a homogeneous (4x4 for 3D-space) Translation matrix that moves a point by Tx units in the x-direction, Ty units in the y-direction,
// and Tz units in the z-direction
func Translate(translation vec3.T) T {
	return T{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, translation.X, translation.Y, translation.Z, 1}
}

// Scale creates a homogeneous 3D scaling matrix
func Scale(scale vec3.T) T {
	return T{scale.X, 0, 0, 0, 0, scale.Y, 0, 0, 0, 0, scale.Z, 0, 0, 0, 0, 1}
}

// Rotate creates a homogenous 3D rotation matrix from euler angles in degrees
func Rotate(rotation vec3.T) T {
	rad := rotation.Scaled(math.Pi / 180.0) // translate rotaiton to radians
	rot := mgl.AnglesToQuat(rad.Z, rad.Y, rad.X, mgl.ZYX).Mat4()
	return T(rot)
}

// Transform creates a new transformation matrix from a position vector, an euler rotation vector, and a scale vector.
func Transform(position, rotation, scale vec3.T) T {
	T := Translate(position)
	R := Rotate(rotation)
	S := Scale(scale)
	// T * R * S
	tr := T.Mul(&R)
	return tr.Mul(&S)
}
