package render

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

/* Represents a shader part of a GLSL program. */
type Shader struct {
    Id          uint32
    stype       ShaderType
    compiled    bool
}

func CreateShader(shaderType ShaderType) *Shader {
    id := gl.CreateShader(uint32(shaderType))
    return &Shader {
        Id:         id,
        stype:      shaderType,
        compiled:   false,
    }
}

/* Compiles and returns a vertex shader from the given source file
   Panics on compilation errors */
func VertexShader(path string) *Shader {
    s := CreateShader(VertexShaderType)
    err := s.CompileFile(path)
    if err != nil {
        panic(err)
    }
    return s
}

/* Compiles and returns a fragment shader from the given source file. 
   Panics on compilation errors */
func FragmentShader(path string) *Shader {
    s := CreateShader(FragmentShaderType)
    err := s.CompileFile(path)
    if err != nil {
        panic(err)
    }
    return s
}

/* Loads and compiles source code from the given file path */
func (shader *Shader) CompileFile(path string) error {
    source, err := ioutil.ReadFile(util.ExePath + path)
    if err != nil {
        return err
    }
    return shader.Compile(string(source))
}

/* Compiles a shader from a source string */
func (shader *Shader) Compile(source string) error {
	csource, free := util.GLString(source)
	gl.ShaderSource(shader.Id, 1, csource, nil)
	gl.CompileShader(shader.Id)
    free()

    /* Check compilation status */
	var status int32
	gl.GetShaderiv(shader.Id, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader.Id, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader.Id, logLength, nil, gl.Str(log))

        return fmt.Errorf("Shader compilation failed.\n** Source: **\n%v\n** Log: **\n%v\n", source, log)
	}

    return nil
}
