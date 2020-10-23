package render

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// ErrUnknownType is returend when an illegal GL type name is used
var ErrUnknownType = fmt.Errorf("Unknown GL Type")

// GLType holds OpenGL type constants
type GLType uint32

// GL Type Constants
const (
	Int8      = GLType(gl.BYTE)
	UInt8     = GLType(gl.UNSIGNED_BYTE)
	Int16     = GLType(gl.SHORT)
	UInt16    = GLType(gl.UNSIGNED_SHORT)
	Int32     = GLType(gl.INT)
	UInt32    = GLType(gl.UNSIGNED_INT)
	Float     = GLType(gl.FLOAT)
	Double    = GLType(gl.DOUBLE)
	Vec2f     = GLType(gl.FLOAT_VEC2)
	Vec3f     = GLType(gl.FLOAT_VEC3)
	Vec4f     = GLType(gl.FLOAT_VEC4)
	Mat3f     = GLType(gl.FLOAT_MAT3)
	Mat4f     = GLType(gl.FLOAT_MAT4)
	Texture2D = GLType(gl.SAMPLER_2D)
)

// Size returns the byte size of the GL type
func (t GLType) Size() int {
	switch t {
	case Int8:
		return 1
	case UInt8:
		return 1
	case Int16:
		return 2
	case UInt16:
		return 2
	case Int32:
		return 4
	case UInt32:
		return 4
	case Float:
		return 4
	case Double:
		return 8
	}
	panic(ErrUnknownType)
}

func (t GLType) String() string {
	switch t {
	case Int8:
		return "int8"
	case UInt8:
		return "uint8"
	case Int16:
		return "int16"
	case UInt16:
		return "uint16"
	case Int32:
		return "int32"
	case UInt32:
		return "uint32"
	case Float:
		return "float"
	case Double:
		return "double"
	case Vec2f:
		return "vec2f"
	case Vec3f:
		return "vec3f"
	case Vec4f:
		return "vec4f"
	case Mat3f:
		return "mat3f"
	case Mat4f:
		return "mat4f"
	case Texture2D:
		return "tex2d"
	default:
		return fmt.Sprintf("unknown:%d", t)
	}
}

// Integer returns the if the type is an integer type
func (t GLType) Integer() bool {
	switch t {
	case Float:
		return false
	case Vec2f:
		return false
	case Vec3f:
		return false
	case Vec4f:
		return false
	case Double:
		return false
	default:
		return true
	}
}

// GLTypeFromString returns the GL identifier & size of a data type name
func GLTypeFromString(name string) (GLType, error) {
	switch name {
	case "byte":
		fallthrough
	case "int8":
		return Int8, nil

	case "ubyte":
		fallthrough
	case "uint8":
		return UInt8, nil

	case "short":
		fallthrough
	case "int16":
		return Int16, nil

	case "ushort":
		fallthrough
	case "uint16":
		return UInt16, nil

	case "int":
		fallthrough
	case "int32":
		return Int32, nil

	case "uint":
		fallthrough
	case "uint32":
		return UInt32, nil

	case "float":
		fallthrough
	case "float32":
		return Float, nil

	case "float64":
		fallthrough
	case "double":
		return Double, nil
	}
	return GLType(0), ErrUnknownType
}
