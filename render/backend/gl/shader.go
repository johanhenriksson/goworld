package gl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/johanhenriksson/goworld/render/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var ErrCompileFailed = errors.New("shader compilation failed")

func CreateShader(shaderType uint32) shader.StageID {
	id := gl.CreateShader(shaderType)
	return shader.StageID(id)
}

func AttachShader(shaderID shader.ShaderID, stageID shader.StageID) {
	gl.AttachShader(uint32(shaderID), uint32(stageID))
}

func CompileShader(id shader.StageID, source, path string) error {
	csource, free := String(source)
	gl.ShaderSource(uint32(id), 1, csource, nil)
	gl.CompileShader(uint32(id))
	free()

	// check compilation status
	var status int32
	gl.GetShaderiv(uint32(id), gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(uint32(id), gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(uint32(id), logLength, nil, gl.Str(log))

		return fmt.Errorf("%w\n** Source: %s**\n%v\n** Log: **\n%v", ErrCompileFailed, path, source, log)
	}

	return nil
}
