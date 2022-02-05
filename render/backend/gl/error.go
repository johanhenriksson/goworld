package gl

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var ErrInvalidEnum = Error{"invalid enum", gl.INVALID_ENUM}
var ErrInvalidValue = Error{"invalid value", gl.INVALID_VALUE}
var ErrInvalidOperation = Error{"invalid operation", gl.INVALID_OPERATION}

type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string {
	return e.Message
}

func GetError() error {
	code := int(gl.GetError())
	if code == 0 {
		return nil
	}

	switch code {
	case gl.INVALID_ENUM:
		return ErrInvalidEnum

	case gl.INVALID_VALUE:
		return ErrInvalidValue

	case gl.INVALID_OPERATION:
		return ErrInvalidOperation

	default:
		return Error{
			Code:    code,
			Message: fmt.Sprintf("OpenGL returned error %d", code),
		}
	}
}
