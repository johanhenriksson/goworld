package gl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/johanhenriksson/goworld/render/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var ErrLinkFailed = errors.New("failed to link program")

func CreateProgram() shader.ShaderID {
	id := gl.CreateProgram()
	if err := GetError(); err != nil {
		panic(err)
	}
	return shader.ShaderID(id)
}

func LinkProgram(id shader.ShaderID) error {
	gl.LinkProgram(uint32(id))
	if err := GetError(); err != nil {
		return err
	}

	// read status
	var status int32
	gl.GetProgramiv(uint32(id), gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(uint32(id), gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(uint32(id), logLength, nil, gl.Str(log))

		return fmt.Errorf("%w: %v", ErrLinkFailed, log)
	}

	return nil
}

func UseProgram(id shader.ShaderID) error {
	gl.UseProgram(uint32(id))
	return GetError()
}

func BindFragDataLocation(id shader.ShaderID, variableName string) error {
	cstr, free := String(variableName)
	defer free()
	gl.BindFragDataLocation(uint32(id), 0, *cstr)
	if err := gl.GetError(); err != gl.NONE {
		return fmt.Errorf("%w: bind fragment data location failed with error %d", shader.ErrUpdateUniform, err)
	}
	return nil
}
