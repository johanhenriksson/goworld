package physics

import (
    "github.com/johanhenriksson/ode"
    mgl "github.com/go-gl/mathgl/mgl32"
    "math"
)

/* ODE Compability Layer Utilities */

const (
    deg2rad = math.Pi / 180.0
    rad2deg = 1.0 / deg2rad
)

/* Convert ODE Vector3 (64 bit) to mgl Vec3 (32 bit) */
func FromOdeVec3(vec ode.Vector3) mgl.Vec3 {
    return mgl.Vec3 {
        float32(vec[0]),
        float32(vec[1]),
        float32(vec[2]),
    }
}

/* Convert mgl Vec3 (32 bit) to an ODE Vector3 (64 bit) */
func ToOdeVec3(vec mgl.Vec3) ode.Vector3 {
    return ode.Vector3 {
        float64(vec[0]),
        float64(vec[1]),
        float64(vec[2]),
    }
}

/* Decompose a 3x3 rotation matrix into euler angles (degrees)
 * http://nghiaho.com/?page_id=846 */
func FromOdeRotation(mat3 ode.Matrix3) mgl.Vec3 {
    x := math.Atan2(mat3[2][1], mat3[2][2])
    y := math.Atan2(-mat3[2][0], math.Sqrt(math.Pow(mat3[2][1], 2) + math.Pow(mat3[2][2], 2)))
    z := math.Atan2(mat3[1][0], mat3[0][0])

    return mgl.Vec3 {
        float32(x * rad2deg),
        float32(y * rad2deg),
        float32(z * rad2deg),
    }
}

func odeV3(x, y, z float32) ode.Vector3 {
    return ode.V3(float64(x), float64(y), float64(z))
}

func odeV4(x, y, z, w float32) ode.Vector4 {
    return ode.V4(float64(x), float64(y), float64(z), float64(w))
}
