package render

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// ErrUnkwownType is returend when an illegal GL type name is used
var ErrUnknownType = fmt.Errorf("Unknown GL Type")

type GLType uint32

const (
	Int8   = GLType(gl.BYTE)
	UInt8  = GLType(gl.UNSIGNED_BYTE)
	Int16  = GLType(gl.SHORT)
	UInt16 = GLType(gl.UNSIGNED_SHORT)
	Int32  = GLType(gl.INT)
	UInt32 = GLType(gl.UNSIGNED_INT)
	Float  = GLType(gl.FLOAT)
	Double = GLType(gl.DOUBLE)
)

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
		fallthrough
	case "unsigned byte":
		return UInt8, nil

	case "short":
		fallthrough
	case "int16":
		return Int16, nil

	case "ushort":
		fallthrough
	case "uint16":
		fallthrough
	case "unsigned short":
		return UInt16, nil

	case "int":
		fallthrough
	case "int32":
		fallthrough
	case "integer":
		return Int32, nil

	case "uint":
		fallthrough
	case "uint32":
		fallthrough
	case "unsigned integer":
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
