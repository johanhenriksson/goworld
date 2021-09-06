package mat4

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"golang.org/x/image/math/f32"

	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
)

// T holds a 4x4 float32 matrix
type T f32.Mat4

// Add performs an element-wise addition of two matrices, this is
// equivalent to iterating over every element of m and adding the corresponding value of m2.
func (m T) Add(m2 T) T {
	return T{
		m[0] + m2[0], m[1] + m2[1], m[2] + m2[2], m[3] + m2[3],
		m[4] + m2[4], m[5] + m2[5], m[6] + m2[6], m[7] + m2[7],
		m[8] + m2[8], m[9] + m2[9], m[10] + m2[10], m[11] + m2[11],
		m[12] + m2[12], m[13] + m2[13], m[14] + m2[14], m[15] + m2[15],
	}
}

// Sub performs an element-wise subtraction of two matrices, this is
// equivalent to iterating over every element of m and subtracting the corresponding value of m2.
func (m T) Sub(m2 T) T {
	return T{
		m[0] - m2[0], m[1] - m2[1], m[2] - m2[2], m[3] - m2[3],
		m[4] - m2[4], m[5] - m2[5], m[6] - m2[6], m[7] - m2[7],
		m[8] - m2[8], m[9] - m2[9], m[10] - m2[10], m[11] - m2[11],
		m[12] - m2[12], m[13] - m2[13], m[14] - m2[14], m[15] - m2[15],
	}
}

// Scale performs a scalar multiplcation of the matrix. This is equivalent to iterating
// over every element of the matrix and multiply it by c.
func (m T) Scale(c float32) T {
	return T{
		m[0] * c, m[1] * c, m[2] * c, m[3] * c,
		m[4] * c, m[5] * c, m[6] * c, m[7] * c,
		m[8] * c, m[9] * c, m[10] * c, m[11] * c,
		m[12] * c, m[13] * c, m[14] * c, m[15] * c,
	}
}

// VMul multiplies a vec4 with the matrix
func (m *T) VMul(v vec4.T) vec4.T {
	return vec4.T{
		X: m[0]*v.X + m[4]*v.Y + m[8]*v.Z + m[12]*v.W,
		Y: m[1]*v.X + m[5]*v.Y + m[9]*v.Z + m[13]*v.W,
		Z: m[2]*v.X + m[6]*v.Y + m[10]*v.Z + m[14]*v.W,
		W: m[3]*v.X + m[7]*v.Y + m[11]*v.Z + m[15]*v.W,
	}
}

// TransformPoint transforms a point to world space
func (m *T) TransformPoint(v vec3.T) vec3.T {
	p := vec4.Extend(v, 1)
	vt := m.VMul(p)
	return vt.XYZ().Scaled(1 / vt.W)
}

// TransformDir transforms a direction vector to world space
func (m *T) TransformDir(v vec3.T) vec3.T {
	p := vec4.Extend(v, 0)
	vt := m.VMul(p)
	return vt.XYZ()
}

// Mul performs a "matrix product" between this matrix and another of the same dimension
func (m *T) Mul(m2 *T) T {
	return T{
		m[0]*m2[0] + m[4]*m2[1] + m[8]*m2[2] + m[12]*m2[3],
		m[1]*m2[0] + m[5]*m2[1] + m[9]*m2[2] + m[13]*m2[3],
		m[2]*m2[0] + m[6]*m2[1] + m[10]*m2[2] + m[14]*m2[3],
		m[3]*m2[0] + m[7]*m2[1] + m[11]*m2[2] + m[15]*m2[3],

		m[0]*m2[4] + m[4]*m2[5] + m[8]*m2[6] + m[12]*m2[7],
		m[1]*m2[4] + m[5]*m2[5] + m[9]*m2[6] + m[13]*m2[7],
		m[2]*m2[4] + m[6]*m2[5] + m[10]*m2[6] + m[14]*m2[7],
		m[3]*m2[4] + m[7]*m2[5] + m[11]*m2[6] + m[15]*m2[7],

		m[0]*m2[8] + m[4]*m2[9] + m[8]*m2[10] + m[12]*m2[11],
		m[1]*m2[8] + m[5]*m2[9] + m[9]*m2[10] + m[13]*m2[11],
		m[2]*m2[8] + m[6]*m2[9] + m[10]*m2[10] + m[14]*m2[11],
		m[3]*m2[8] + m[7]*m2[9] + m[11]*m2[10] + m[15]*m2[11],

		m[0]*m2[12] + m[4]*m2[13] + m[8]*m2[14] + m[12]*m2[15],
		m[1]*m2[12] + m[5]*m2[13] + m[9]*m2[14] + m[13]*m2[15],
		m[2]*m2[12] + m[6]*m2[13] + m[10]*m2[14] + m[14]*m2[15],
		m[3]*m2[12] + m[7]*m2[13] + m[11]*m2[14] + m[15]*m2[15],
	}
}

// Transpose produces the transpose of this matrix. For any MxN matrix
// the transpose is an NxM matrix with the rows swapped with the columns. For instance
// the transpose of the Mat3x2 is a Mat2x3 like so:
//
//    [[a b]]    [[a c e]]
//    [[c d]] =  [[b d f]]
//    [[e f]]
func (m *T) Transpose() T {
	return T{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

// Det returns the determinant of a matrix. It is a measure of a square matrix's
// singularity and invertability, among other things. In this library, the
// determinant is hard coded based on pre-computed cofactor expansion, and uses
// no loops. Of course, the addition and multiplication must still be done.
func (m *T) Det() float32 {
	return m[0]*m[5]*m[10]*m[15] - m[0]*m[5]*m[11]*m[14] - m[0]*m[6]*m[9]*m[15] + m[0]*m[6]*m[11]*m[13] +
		m[0]*m[7]*m[9]*m[14] - m[0]*m[7]*m[10]*m[13] - m[1]*m[4]*m[10]*m[15] + m[1]*m[4]*m[11]*m[14] +
		m[1]*m[6]*m[8]*m[15] - m[1]*m[6]*m[11]*m[12] - m[1]*m[7]*m[8]*m[14] + m[1]*m[7]*m[10]*m[12] +
		m[2]*m[4]*m[9]*m[15] - m[2]*m[4]*m[11]*m[13] - m[2]*m[5]*m[8]*m[15] + m[2]*m[5]*m[11]*m[12] +
		m[2]*m[7]*m[8]*m[13] - m[2]*m[7]*m[9]*m[12] - m[3]*m[4]*m[9]*m[14] + m[3]*m[4]*m[10]*m[13] +
		m[3]*m[5]*m[8]*m[14] - m[3]*m[5]*m[10]*m[12] - m[3]*m[6]*m[8]*m[13] + m[3]*m[6]*m[9]*m[12]
}

// Invert computes the inverse of a square matrix. An inverse is a square matrix such that when multiplied by the
// original, yields the identity.
//
// M_inv * M = M * M_inv = I
//
// In this library, the math is precomputed, and uses no loops, though the multiplications, additions, determinant calculation, and scaling
// are still done. This can still be (relatively) expensive for a 4x4.
//
// This function checks the determinant to see if the matrix is invertible.
// If the determinant is 0.0, this function returns the zero matrix. However, due to floating point errors, it is
// entirely plausible to get a false positive or negative.
// In the future, an alternate function may be written which takes in a pre-computed determinant.
func (m *T) Invert() T {
	det := m.Det()
	if math.Equal(det, float32(0.0)) {
		return T{}
	}

	retMat := T{
		-m[7]*m[10]*m[13] + m[6]*m[11]*m[13] + m[7]*m[9]*m[14] - m[5]*m[11]*m[14] - m[6]*m[9]*m[15] + m[5]*m[10]*m[15],
		m[3]*m[10]*m[13] - m[2]*m[11]*m[13] - m[3]*m[9]*m[14] + m[1]*m[11]*m[14] + m[2]*m[9]*m[15] - m[1]*m[10]*m[15],
		-m[3]*m[6]*m[13] + m[2]*m[7]*m[13] + m[3]*m[5]*m[14] - m[1]*m[7]*m[14] - m[2]*m[5]*m[15] + m[1]*m[6]*m[15],
		m[3]*m[6]*m[9] - m[2]*m[7]*m[9] - m[3]*m[5]*m[10] + m[1]*m[7]*m[10] + m[2]*m[5]*m[11] - m[1]*m[6]*m[11],
		m[7]*m[10]*m[12] - m[6]*m[11]*m[12] - m[7]*m[8]*m[14] + m[4]*m[11]*m[14] + m[6]*m[8]*m[15] - m[4]*m[10]*m[15],
		-m[3]*m[10]*m[12] + m[2]*m[11]*m[12] + m[3]*m[8]*m[14] - m[0]*m[11]*m[14] - m[2]*m[8]*m[15] + m[0]*m[10]*m[15],
		m[3]*m[6]*m[12] - m[2]*m[7]*m[12] - m[3]*m[4]*m[14] + m[0]*m[7]*m[14] + m[2]*m[4]*m[15] - m[0]*m[6]*m[15],
		-m[3]*m[6]*m[8] + m[2]*m[7]*m[8] + m[3]*m[4]*m[10] - m[0]*m[7]*m[10] - m[2]*m[4]*m[11] + m[0]*m[6]*m[11],
		-m[7]*m[9]*m[12] + m[5]*m[11]*m[12] + m[7]*m[8]*m[13] - m[4]*m[11]*m[13] - m[5]*m[8]*m[15] + m[4]*m[9]*m[15],
		m[3]*m[9]*m[12] - m[1]*m[11]*m[12] - m[3]*m[8]*m[13] + m[0]*m[11]*m[13] + m[1]*m[8]*m[15] - m[0]*m[9]*m[15],
		-m[3]*m[5]*m[12] + m[1]*m[7]*m[12] + m[3]*m[4]*m[13] - m[0]*m[7]*m[13] - m[1]*m[4]*m[15] + m[0]*m[5]*m[15],
		m[3]*m[5]*m[8] - m[1]*m[7]*m[8] - m[3]*m[4]*m[9] + m[0]*m[7]*m[9] + m[1]*m[4]*m[11] - m[0]*m[5]*m[11],
		m[6]*m[9]*m[12] - m[5]*m[10]*m[12] - m[6]*m[8]*m[13] + m[4]*m[10]*m[13] + m[5]*m[8]*m[14] - m[4]*m[9]*m[14],
		-m[2]*m[9]*m[12] + m[1]*m[10]*m[12] + m[2]*m[8]*m[13] - m[0]*m[10]*m[13] - m[1]*m[8]*m[14] + m[0]*m[9]*m[14],
		m[2]*m[5]*m[12] - m[1]*m[6]*m[12] - m[2]*m[4]*m[13] + m[0]*m[6]*m[13] + m[1]*m[4]*m[14] - m[0]*m[5]*m[14],
		-m[2]*m[5]*m[8] + m[1]*m[6]*m[8] + m[2]*m[4]*m[9] - m[0]*m[6]*m[9] - m[1]*m[4]*m[10] + m[0]*m[5]*m[10],
	}

	return retMat.Scale(1 / det)
}

// ApproxEqual performs an element-wise approximate equality test between two matrices,
// as if FloatEqual had been used.
func (m *T) ApproxEqual(m2 *T) bool {
	for i := range m {
		if !math.Equal(m[i], m2[i]) {
			return false
		}
	}
	return true
}

// ApproxEqualThreshold performs an element-wise approximate equality test between two matrices
// with a given epsilon threshold, as if FloatEqualThreshold had been used.
func (m *T) ApproxEqualThreshold(m2 *T, threshold float32) bool {
	for i := range m {
		if !math.EqualThreshold(m[i], m2[i], threshold) {
			return false
		}
	}
	return true
}

// At returns the matrix element at the given row and column.
// This is equivalent to mat[col * numRow + row] where numRow is constant
// (E.G. for a Mat3x2 it's equal to 3)
//
// This method is garbage-in garbage-out. For instance, on a T asking for
// At(5,0) will work just like At(1,1). Or it may panic if it's out of bounds.
func (m *T) At(row, col int) float32 {
	return m[col*4+row]
}

// Set sets the corresponding matrix element at the given row and column.
func (m *T) Set(row, col int, value float32) {
	m[col*4+row] = value
}

// Index returns the index of the given row and column, to be used with direct
// access. E.G. Index(0,0) = 0.
func (m *T) Index(row, col int) int {
	return col*4 + row
}

// Row returns a vector representing the corresponding row (starting at row 0).
// This package makes no distinction between row and column vectors, so it
// will be a normal VecM for a MxN matrix.
func (m *T) Row(row int) vec4.T {
	return vec4.T{
		X: m[row+0],
		Y: m[row+4],
		Z: m[row+8],
		W: m[row+12],
	}
}

// Rows decomposes a matrix into its corresponding row vectors.
// This is equivalent to calling mat.Row for each row.
func (m *T) Rows() (row0, row1, row2, row3 vec4.T) {
	return m.Row(0), m.Row(1), m.Row(2), m.Row(3)
}

// Col returns a vector representing the corresponding column (starting at col 0).
// This package makes no distinction between row and column vectors, so it
// will be a normal VecN for a MxN matrix.
func (m *T) Col(col int) vec4.T {
	return vec4.T{
		X: m[col*4+0],
		Y: m[col*4+1],
		Z: m[col*4+2],
		W: m[col*4+3],
	}
}

// Cols decomposes a matrix into its corresponding column vectors.
// This is equivalent to calling mat.Col for each column.
func (m *T) Cols() (col0, col1, col2, col3 vec4.T) {
	return m.Col(0), m.Col(1), m.Col(2), m.Col(3)
}

// Trace is a basic operation on a square matrix that simply
// sums up all elements on the main diagonal (meaning all elements such that row==col).
func (m *T) Trace() float32 {
	return m[0] + m[5] + m[10] + m[15]
}

// Abs returns the element-wise absolute value of this matrix
func (m *T) Abs() T {
	return T{
		math.Abs(m[0]), math.Abs(m[1]), math.Abs(m[2]), math.Abs(m[3]),
		math.Abs(m[4]), math.Abs(m[5]), math.Abs(m[6]), math.Abs(m[7]),
		math.Abs(m[8]), math.Abs(m[9]), math.Abs(m[10]), math.Abs(m[11]),
		math.Abs(m[12]), math.Abs(m[13]), math.Abs(m[14]), math.Abs(m[15]),
	}
}

// String pretty prints the matrix
func (m T) String() string {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 4, 4, 1, ' ', tabwriter.AlignRight)
	for i := 0; i < 4; i++ {
		r := m.Row(i)
		fmt.Fprintf(w, "%f\t", r.X)
		fmt.Fprintf(w, "%f\t", r.Y)
		fmt.Fprintf(w, "%f\t", r.Z)
		fmt.Fprintf(w, "%f\t", r.W)
	}
	w.Flush()
	return buf.String()
}

// Right extracts the right vector from a transformation matrix
func (m *T) Right() vec3.T {
	return vec3.T{
		X: m[4*0+0],
		Y: m[4*1+0],
		Z: m[4*2+0],
	}
}

// Up extracts the up vector from a transformation matrix
func (m *T) Up() vec3.T {
	return vec3.T{
		X: m[4*0+1],
		Y: m[4*1+1],
		Z: m[4*2+1],
	}
}

// Forward extracts the forward vector from a transformation matrix
func (m *T) Forward() vec3.T {
	return vec3.T{
		X: -m[4*0+2],
		Y: -m[4*1+2],
		Z: -m[4*2+2],
	}
}
