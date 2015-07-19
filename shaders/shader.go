package shaders

import (
    "fmt"
    "strings"
    "io/ioutil"

	"github.com/go-gl/gl/v4.1-core/gl"

    "github.com/johanhenriksson/goworld/util"
)

type ShaderType uint32

const VertexShaderType   ShaderType = gl.VERTEX_SHADER
const FragmentShaderType ShaderType = gl.FRAGMENT_SHADER

type Shader struct {
    Id          uint32
    stype       ShaderType
    compiled    bool
}

func Create(shaderType ShaderType) *Shader {
    id := gl.CreateShader(uint32(shaderType))
    return &Shader {
        Id:         id,
        stype:      shaderType,
        compiled:   false,
    }
}

func VertexShader(path string) *Shader {
    s := Create(VertexShaderType)
    err := s.CompileFile(path)
    if err != nil {
        panic(err)
    }
    return s
}

func FragmentShader(path string) *Shader {
    s := Create(FragmentShaderType)
    err := s.CompileFile(path)
    if err != nil {
        panic(err)
    }
    return s
}


func (shader *Shader) CompileFile(path string) error {
    source, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    return shader.Compile(string(source))
}

func (shader *Shader) Compile(source string) error {
	csource := util.GLString(source)
	gl.ShaderSource(shader.Id, 1, &csource, nil)
	gl.CompileShader(shader.Id)

	var status int32
	gl.GetShaderiv(shader.Id, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader.Id, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader.Id, logLength, nil, gl.Str(log))

		return fmt.Errorf("Failed to compile %v: %v", source, log)
	}

    return nil
}
