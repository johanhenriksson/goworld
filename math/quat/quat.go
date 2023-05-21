// Based on code from github.com/go-gl/mathgl:
// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package quat

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// RotationOrder is the order in which rotations will be transformed for the
// purposes of AnglesToQuat.
type RotationOrder int

// The RotationOrder constants represent a series of rotations along the given
// axes for the use of AnglesToQuat.
const (
	XYX RotationOrder = iota
	XYZ
	XZX
	XZY
	YXY
	YXZ
	YZY
	YZX
	ZYZ
	ZYX
	ZXZ
	ZXY
)

// T represents a Quaternion, which is an extension of the imaginary numbers;
// there's all sorts of interesting theory behind it. In 3D graphics we mostly
// use it as a cheap way of representing rotation since quaternions are cheaper
// to multiply by, and easier to interpolate than matrices.
//
// A Quaternion has two parts: W, the so-called scalar component, and "V", the
// vector component. The vector component is considered to be the part in 3D
// space, while W (loosely interpreted) is its 4D coordinate.
type T struct {
	W float32
	V vec3.T
}

// Ident returns the quaternion identity: W=1; V=(0,0,0).
//
// As with all identities, multiplying any quaternion by this will yield the same
// quaternion you started with.
func Ident() T {
	return T{1., vec3.New(0, 0, 0)}
}

// Rotate creates an angle from an axis and an angle relative to that axis.
//
// This is cheaper than HomogRotate3D.
func Rotate(angle float32, axis vec3.T) T {
	// angle = (float32(math.Pi) * angle) / 180.0

	c, s := math.Cos(angle/2), math.Sin(angle/2)

	return T{c, axis.Scaled(s)}
}

// X is a convenient alias for q.V[0]
func (q T) X() float32 {
	return q.V.X
}

// Y is a convenient alias for q.V[1]
func (q T) Y() float32 {
	return q.V.Y
}

// Z is a convenient alias for q.V[2]
func (q T) Z() float32 {
	return q.V.X
}

// Add adds two quaternions. It's no more complicated than
// adding their W and V components.
func (q1 T) Add(q2 T) T {
	return T{q1.W + q2.W, q1.V.Add(q2.V)}
}

// Sub subtracts two quaternions. It's no more complicated than
// subtracting their W and V components.
func (q1 T) Sub(q2 T) T {
	return T{q1.W - q2.W, q1.V.Sub(q2.V)}
}

// Mul multiplies two quaternions. This can be seen as a rotation. Note that
// Multiplication is NOT commutative, meaning q1.Mul(q2) does not necessarily
// equal q2.Mul(q1).
func (q1 T) Mul(q2 T) T {
	return T{q1.W*q2.W - vec3.Dot(q1.V, q2.V), vec3.Cross(q1.V, q2.V).Add(q2.V.Scaled(q1.W)).Add(q1.V.Scaled(q2.W))}
}

// Scale every element of the quaternion by some constant factor.
func (q1 T) Scale(c float32) T {
	return T{q1.W * c, vec3.New(q1.V.X*c, q1.V.Y*c, q1.V.Z*c)}
}

// Conjugate returns the conjugate of a quaternion. Equivalent to
// Quat{q1.W, q1.V.Mul(-1)}.
func (q1 T) Conjugate() T {
	return T{q1.W, q1.V.Scaled(-1)}
}

// Len gives the Length of the quaternion, also known as its Norm. This is the
// same thing as the Len of a Vec4.
func (q1 T) Len() float32 {
	return math.Sqrt(q1.W*q1.W + vec3.Dot(q1.V, q1.V))
}

// Norm is an alias for Len() since both are very common terms.
func (q1 T) Norm() float32 {
	return q1.Len()
}

// Normalize the quaternion, returning its versor (unit quaternion).
//
// This is the same as normalizing it as a Vec4.
func (q1 T) Normalize() T {
	length := q1.Len()

	if math.Equal(1, length) {
		return q1
	}
	if length == 0 {
		return Ident()
	}
	if length == math.InfPos {
		length = math.MaxValue
	}

	return T{q1.W * 1 / length, q1.V.Scaled(1 / length)}
}

// Inverse of a quaternion. The inverse is equivalent
// to the conjugate divided by the square of the length.
//
// This method computes the square norm by directly adding the sum
// of the squares of all terms instead of actually squaring q1.Len(),
// both for performance and precision.
func (q1 T) Inverse() T {
	return q1.Conjugate().Scale(1 / q1.Dot(q1))
}

// Rotate a vector by the rotation this quaternion represents.
// This will result in a 3D vector. Strictly speaking, this is
// equivalent to q1.v.q* where the "."" is quaternion multiplication and v is interpreted
// as a quaternion with W 0 and V v. In code:
// q1.Mul(Quat{0,v}).Mul(q1.Conjugate()), and
// then retrieving the imaginary (vector) part.
//
// In practice, we hand-compute this in the general case and simplify
// to save a few operations.
func (q1 T) Rotate(v vec3.T) vec3.T {
	cross := vec3.Cross(q1.V, v)
	// v + 2q_w * (q_v x v) + 2q_v x (q_v x v)
	return v.Add(cross.Scaled(2 * q1.W)).Add(vec3.Cross(q1.V.Scaled(2), cross))
}

// Mat4 returns the homogeneous 3D rotation matrix corresponding to the
// quaternion.
func (q1 T) Mat4() mat4.T {
	w, x, y, z := q1.W, q1.V.X, q1.V.Y, q1.V.Z
	return mat4.T{
		1 - 2*y*y - 2*z*z, 2*x*y + 2*w*z, 2*x*z - 2*w*y, 0,
		2*x*y - 2*w*z, 1 - 2*x*x - 2*z*z, 2*y*z + 2*w*x, 0,
		2*x*z + 2*w*y, 2*y*z - 2*w*x, 1 - 2*x*x - 2*y*y, 0,
		0, 0, 0, 1,
	}
}

// Dot product between two quaternions, equivalent to if this was a Vec4.
func (q1 T) Dot(q2 T) float32 {
	return q1.W*q2.W + vec3.Dot(q1.V, q2.V)
}

// ApproxEqual returns whether the quaternions are approximately equal, as if
// FloatEqual was called on each matching element
func (q1 T) ApproxEqual(q2 T) bool {
	return math.Equal(q1.W, q2.W) && q1.V.ApproxEqual(q2.V)
}

// OrientationEqual returns whether the quaternions represents the same orientation
//
// Different values can represent the same orientation (q == -q) because quaternions avoid singularities
// and discontinuities involved with rotation in 3 dimensions by adding extra dimensions
func (q1 T) OrientationEqual(q2 T) bool {
	return q1.OrientationEqualThreshold(q2, math.Epsilon)
}

// OrientationEqualThreshold returns whether the quaternions represents the same orientation with a given tolerence
func (q1 T) OrientationEqualThreshold(q2 T, epsilon float32) bool {
	return math.Abs(q1.Normalize().Dot(q2.Normalize())) > 1-math.Epsilon
}

// Slerp is *S*pherical *L*inear Int*erp*olation, a method of interpolating
// between two quaternions. This always takes the straightest path on the sphere between
// the two quaternions, and maintains constant velocity.
//
// However, it's expensive and Slerp(q1,q2) is not the same as Slerp(q2,q1)
func Slerp(q1, q2 T, amount float32) T {
	q1, q2 = q1.Normalize(), q2.Normalize()
	dot := q1.Dot(q2)

	// If the inputs are too close for comfort, linearly interpolate and normalize the result.
	if dot > 0.9995 {
		return Nlerp(q1, q2, amount)
	}

	// This is here for precision errors, I'm perfectly aware that *technically* the dot is bound [-1,1], but since Acos will freak out if it's not (even if it's just a liiiiitle bit over due to normal error) we need to clamp it
	dot = math.Clamp(dot, -1, 1)

	theta := math.Acos(dot) * amount
	c, s := math.Cos(theta), math.Sin(theta)
	rel := q2.Sub(q1.Scale(dot)).Normalize()

	return q1.Scale(c).Add(rel.Scale(s))
}

// Lerp is a *L*inear Int*erp*olation between two Quaternions, cheap and simple.
//
// Not excessively useful, but uses can be found.
func Lerp(q1, q2 T, amount float32) T {
	return q1.Add(q2.Sub(q1).Scale(amount))
}

// Nlerp is a *Normalized* *L*inear Int*erp*olation between two Quaternions. Cheaper than Slerp
// and usually just as good. This is literally Lerp with Normalize() called on it.
//
// Unlike Slerp, constant velocity isn't maintained, but it's much faster and
// Nlerp(q1,q2) and Nlerp(q2,q1) return the same path. You should probably
// use this more often unless you're suffering from choppiness due to the
// non-constant velocity problem.
func Nlerp(q1, q2 T, amount float32) T {
	return Lerp(q1, q2, amount).Normalize()
}

// FromAngles performs a rotation in the specified order. If the order is not
// a valid RotationOrder, this function will panic
//
// The rotation "order" is more of an axis descriptor. For instance XZX would
// tell the function to interpret angle1 as a rotation about the X axis, angle2 about
// the Z axis, and angle3 about the X axis again.
//
// Based off the code for the Matlab function "angle2quat", though this implementation
// only supports 3 single angles as opposed to multiple angles.
func FromAngles(angle1, angle2, angle3 float32, order RotationOrder) T {
	var s [3]float32
	var c [3]float32

	s[0], c[0] = math.Sincos(angle1 / 2)
	s[1], c[1] = math.Sincos(angle2 / 2)
	s[2], c[2] = math.Sincos(angle3 / 2)

	ret := T{}
	switch order {
	case ZYX:
		ret.W = c[0]*c[1]*c[2] + s[0]*s[1]*s[2]
		ret.V = vec3.T{
			X: c[0]*c[1]*s[2] - s[0]*s[1]*c[2],
			Y: c[0]*s[1]*c[2] + s[0]*c[1]*s[2],
			Z: s[0]*c[1]*c[2] - c[0]*s[1]*s[2],
		}
	case ZYZ:
		ret.W = c[0]*c[1]*c[2] - s[0]*c[1]*s[2]
		ret.V = vec3.T{
			X: c[0]*s[1]*s[2] - s[0]*s[1]*c[2],
			Y: c[0]*s[1]*c[2] + s[0]*s[1]*s[2],
			Z: s[0]*c[1]*c[2] + c[0]*c[1]*s[2],
		}
	case ZXY:
		ret.W = c[0]*c[1]*c[2] - s[0]*s[1]*s[2]
		ret.V = vec3.T{
			X: c[0]*s[1]*c[2] - s[0]*c[1]*s[2],
			Y: c[0]*c[1]*s[2] + s[0]*s[1]*c[2],
			Z: c[0]*s[1]*s[2] + s[0]*c[1]*c[2],
		}
	case ZXZ:
		ret.W = c[0]*c[1]*c[2] - s[0]*c[1]*s[2]
		ret.V = vec3.T{
			X: c[0]*s[1]*c[2] + s[0]*s[1]*s[2],
			Y: s[0]*s[1]*c[2] - c[0]*s[1]*s[2],
			Z: c[0]*c[1]*s[2] + s[0]*c[1]*c[2],
		}
	case YXZ:
		ret.W = c[0]*c[1]*c[2] + s[0]*s[1]*s[2]
		ret.V = vec3.T{
			X: c[0]*s[1]*c[2] + s[0]*c[1]*s[2],
			Y: s[0]*c[1]*c[2] - c[0]*s[1]*s[2],
			Z: c[0]*c[1]*s[2] - s[0]*s[1]*c[2],
		}
	case YXY:
		ret.W = c[0]*c[1]*c[2] - s[0]*c[1]*s[2]
		ret.V = vec3.T{
			X: c[0]*s[1]*c[2] + s[0]*s[1]*s[2],
			Y: s[0]*c[1]*c[2] + c[0]*c[1]*s[2],
			Z: c[0]*s[1]*s[2] - s[0]*s[1]*c[2],
		}
	case YZX:
		ret.W = c[0]*c[1]*c[2] - s[0]*s[1]*s[2]
		ret.V = vec3.T{
			X: c[0]*c[1]*s[2] + s[0]*s[1]*c[2],
			Y: c[0]*s[1]*s[2] + s[0]*c[1]*c[2],
			Z: c[0]*s[1]*c[2] - s[0]*c[1]*s[2],
		}
	case YZY:
		ret.W = c[0]*c[1]*c[2] - s[0]*c[1]*s[2]
		ret.V = vec3.T{
			X: s[0]*s[1]*c[2] - c[0]*s[1]*s[2],
			Y: c[0]*c[1]*s[2] + s[0]*c[1]*c[2],
			Z: c[0]*s[1]*c[2] + s[0]*s[1]*s[2],
		}
	case XYZ:
		ret.W = c[0]*c[1]*c[2] - s[0]*s[1]*s[2]
		ret.V = vec3.T{
			X: c[0]*s[1]*s[2] + s[0]*c[1]*c[2],
			Y: c[0]*s[1]*c[2] - s[0]*c[1]*s[2],
			Z: c[0]*c[1]*s[2] + s[0]*s[1]*c[2],
		}
	case XYX:
		ret.W = c[0]*c[1]*c[2] - s[0]*c[1]*s[2]
		ret.V = vec3.T{
			X: c[0]*c[1]*s[2] + s[0]*c[1]*c[2],
			Y: c[0]*s[1]*c[2] + s[0]*s[1]*s[2],
			Z: s[0]*s[1]*c[2] - c[0]*s[1]*s[2],
		}
	case XZY:
		ret.W = c[0]*c[1]*c[2] + s[0]*s[1]*s[2]
		ret.V = vec3.T{
			X: s[0]*c[1]*c[2] - c[0]*s[1]*s[2],
			Y: c[0]*c[1]*s[2] - s[0]*s[1]*c[2],
			Z: c[0]*s[1]*c[2] + s[0]*c[1]*s[2],
		}
	case XZX:
		ret.W = c[0]*c[1]*c[2] - s[0]*c[1]*s[2]
		ret.V = vec3.T{
			X: c[0]*c[1]*s[2] + s[0]*c[1]*c[2],
			Y: c[0]*s[1]*s[2] - s[0]*s[1]*c[2],
			Z: c[0]*s[1]*c[2] + s[0]*s[1]*s[2],
		}
	default:
		panic("Unsupported rotation order")
	}
	return ret
}

// FromMat4 converts a pure rotation matrix into a quaternion
func FromMat4(m mat4.T) T {
	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/matrixToQuaternion/index.htm

	if tr := m[0] + m[5] + m[10]; tr > 0 {
		s := 0.5 / math.Sqrt(tr+1.0)
		return T{
			0.25 / s,
			vec3.T{
				X: (m[6] - m[9]) * s,
				Y: (m[8] - m[2]) * s,
				Z: (m[1] - m[4]) * s,
			},
		}
	}

	if (m[0] > m[5]) && (m[0] > m[10]) {
		s := 2.0 * math.Sqrt(1.0+m[0]-m[5]-m[10])
		return T{
			(m[6] - m[9]) / s,
			vec3.T{
				X: 0.25 * s,
				Y: (m[4] + m[1]) / s,
				Z: (m[8] + m[2]) / s,
			},
		}
	}

	if m[5] > m[10] {
		s := 2.0 * math.Sqrt(1.0+m[5]-m[0]-m[10])
		return T{
			(m[8] - m[2]) / s,
			vec3.T{
				X: (m[4] + m[1]) / s,
				Y: 0.25 * s,
				Z: (m[9] + m[6]) / s,
			},
		}

	}

	s := 2.0 * math.Sqrt(1.0+m[10]-m[0]-m[5])
	return T{
		(m[1] - m[4]) / s,
		vec3.T{
			X: (m[8] + m[2]) / s,
			Y: (m[9] + m[6]) / s,
			Z: 0.25 * s,
		},
	}
}

// LookAtV creates a rotation from an eye vector to a center vector
//
// It assumes the front of the rotated object at Z- and up at Y+
func LookAtV(eye, center, up vec3.T) T {
	// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-17-quaternions/#I_need_an_equivalent_of_gluLookAt__How_do_I_orient_an_object_towards_a_point__
	// https://bitbucket.org/sinbad/ogre/src/d2ef494c4a2f5d6e2f0f17d3bfb9fd936d5423bb/OgreMain/src/OgreCamera.cpp?at=default#cl-161

	direction := center.Sub(eye).Normalized()

	// Find the rotation between the front of the object (that we assume towards Z-,
	// but this depends on your model) and the desired direction
	rotDir := BetweenVectors(vec3.UnitZN, direction)

	// Recompute up so that it's perpendicular to the direction
	// You can skip that part if you really want to force up
	//right := direction.Cross(up)
	//up = right.Cross(direction)

	// Because of the 1rst rotation, the up is probably completely screwed up.
	// Find the rotation between the "up" of the rotated object, and the desired up
	upCur := rotDir.Rotate(vec3.Zero)
	rotUp := BetweenVectors(upCur, up)

	rotTarget := rotUp.Mul(rotDir) // remember, in reverse order.
	return rotTarget.Inverse()     // camera rotation should be inversed!
}

// BetweenVectors calculates the rotation between two vectors
func BetweenVectors(start, dest vec3.T) T {
	// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-17-quaternions/#I_need_an_equivalent_of_gluLookAt__How_do_I_orient_an_object_towards_a_point__
	// https://github.com/g-truc/glm/blob/0.9.5/glm/gtx/quaternion.inl#L225
	// https://bitbucket.org/sinbad/ogre/src/d2ef494c4a2f5d6e2f0f17d3bfb9fd936d5423bb/OgreMain/include/OgreVector3.h?at=default#cl-654

	start = start.Normalized()
	dest = dest.Normalized()
	epsilon := float32(0.001)

	cosTheta := vec3.Dot(start, dest)
	if cosTheta < -1.0+epsilon {
		// special case when vectors in opposite directions:
		// there is no "ideal" rotation axis
		// So guess one; any will do as long as it's perpendicular to start
		axis := vec3.Cross(vec3.UnitX, start)
		if vec3.Dot(axis, axis) < epsilon {
			// bad luck, they were parallel, try again!
			axis = vec3.Cross(vec3.UnitY, start)
		}

		return Rotate(math.Pi, axis.Normalized())
	}

	axis := vec3.Cross(start, dest)
	s := float32(math.Sqrt(float32(1.0+cosTheta) * 2.0))

	return T{
		s * 0.5,
		axis.Scaled(1.0 / s),
	}
}

func (q T) ToAngles(order RotationOrder) vec3.T {
	// this function was adapted from a Go port of Three.js math, github.com/tengge1/go-three-math
	// Copyright 2017-2020 The ShadowEditor Authors. All rights reserved.
	// Use of e source code is governed by a MIT-style
	// license that can be found in the LICENSE file.

	// assumes the upper 3x3 of m is a pure rotation matrix (i.e, unscaled)
	te := q.Mat4()
	m11, m12, m13 := te[0], te[4], te[8]
	m21, m22, m23 := te[1], te[5], te[9]
	m31, m32, m33 := te[2], te[6], te[10]

	e := vec3.Zero
	switch order {
	default:
		panic("unsupported rotation order")
	case XYZ:
		e.Y = math.Asin(math.Clamp(m13, -1, 1))

		if math.Abs(m13) < 0.9999999 {
			e.X = math.Atan2(-m23, m33)
			e.Z = math.Atan2(-m12, m11)
		} else {
			e.X = math.Atan2(m32, m22)
			e.Z = 0
		}
	case YXZ:
		e.X = math.Asin(-math.Clamp(m23, -1, 1))

		if math.Abs(m23) < 0.9999999 {
			e.Y = math.Atan2(m13, m33)
			e.Z = math.Atan2(m21, m22)
		} else {
			e.Y = math.Atan2(-m31, m11)
			e.Z = 0
		}
	case ZXY:
		e.X = math.Asin(math.Clamp(m32, -1, 1))

		if math.Abs(m32) < 0.9999999 {
			e.Y = math.Atan2(-m31, m33)
			e.Z = math.Atan2(-m12, m22)
		} else {
			e.Y = 0
			e.Z = math.Atan2(m21, m11)
		}
	case ZYX:
		e.Y = math.Asin(-math.Clamp(m31, -1, 1))

		if math.Abs(m31) < 0.9999999 {
			e.X = math.Atan2(m32, m33)
			e.Z = math.Atan2(m21, m11)
		} else {
			e.X = 0
			e.Z = math.Atan2(-m12, m22)
		}
	case YZX:
		e.Z = math.Asin(math.Clamp(m21, -1, 1))

		if math.Abs(m21) < 0.9999999 {
			e.X = math.Atan2(-m23, m22)
			e.Y = math.Atan2(-m31, m11)
		} else {
			e.X = 0
			e.Y = math.Atan2(m13, m33)
		}
	case XZY:
		e.Z = math.Asin(-math.Clamp(m12, -1, 1))

		if math.Abs(m12) < 0.9999999 {
			e.X = math.Atan2(m32, m22)
			e.Y = math.Atan2(m13, m11)
		} else {
			e.X = math.Atan2(-m23, m33)
			e.Y = 0
		}
	}

	return e
}

func (q T) Euler() vec3.T {
	// convert radians to degrees
	return q.ToAngles(ZXY).Scaled(180.0 / math.Pi)
}

func Euler(x, y, z float32) T {
	return FromAngles(math.DegToRad(z), math.DegToRad(x), math.DegToRad(y), ZXY)
}
