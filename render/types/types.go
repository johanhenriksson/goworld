package types

import (
	"errors"
	"fmt"
)

// ErrUnknownType is returend when an illegal GL type name is used
var ErrUnknownType = errors.New("unknown data type")

// Type holds OpenGL type constants
type Type uint32

// GL Type Constants
const (
	_ Type = iota
	Bool
	Int8
	UInt8
	Int16
	UInt16
	Int32
	UInt32
	Float
	Vec2f
	Vec3f
	Vec4f
	Mat3f
	Mat4f
	Double
	Texture2D
)

// Size returns the byte size of the GL type
func (t Type) Size() int {
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
	panic(fmt.Errorf("unknown size for GL type %s", t))
}

func (t Type) String() string {
	switch t {
	case Bool:
		return "bool"
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
func (t Type) Integer() bool {
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

// TypeFromString returns the GL identifier & size of a data type name
func TypeFromString(name string) (Type, error) {
	switch name {
	case "bool":
		return Bool, nil

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
	return Type(0), ErrUnknownType
}
