package render

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/util"
)

// ShaderType indicates the type of shader program
type ShaderType uint32

// VertexShaderType is a Vertex Shader
const VertexShaderType ShaderType = gl.VERTEX_SHADER

// FragmentShaderType is a Fragment Shader
const FragmentShaderType ShaderType = gl.FRAGMENT_SHADER

const GeometryShaderType ShaderType = gl.GEOMETRY_SHADER

// ShaderStage represents a shader part of a GLSL program
type ShaderStage struct {
	ID       uint32
	stype    ShaderType
	compiled bool
}

// CreateShaderStage creates a new empty shader
func CreateShaderStage(shaderType ShaderType) *ShaderStage {
	id := gl.CreateShader(uint32(shaderType))
	return &ShaderStage{
		ID:       id,
		stype:    shaderType,
		compiled: false,
	}
}

// VertexShader compiles and returns a vertex shader from the given source file
// Panics on compilation errors
func VertexShader(path string) *ShaderStage {
	s := CreateShaderStage(VertexShaderType)
	err := s.CompileFile(path)
	if err != nil {
		panic(err)
	}
	return s
}

// FragmentShader compiles and returns a fragment shader from the given source file.
// Panics on compilation errors
func FragmentShader(path string) *ShaderStage {
	s := CreateShaderStage(FragmentShaderType)
	err := s.CompileFile(path)
	if err != nil {
		panic(err)
	}
	return s
}

func GeometryShader(path string) *ShaderStage {
	s := CreateShaderStage(GeometryShaderType)
	err := s.CompileFile(path)
	if err != nil {
		panic(err)
	}
	return s
}

// CompileFile loads and compiles source code from the given file path
func (shader *ShaderStage) CompileFile(path string) error {
	source, err := ioutil.ReadFile(util.ExePath + path)
	if err != nil {
		return err
	}
	return shader.Compile(string(source), path)
}

// Compile a shader from a source string
func (shader *ShaderStage) Compile(source, path string) error {
	csource, free := util.GLString(source)
	gl.ShaderSource(shader.ID, 1, csource, nil)
	gl.CompileShader(shader.ID)
	free()

	/* Check compilation status */
	var status int32
	gl.GetShaderiv(shader.ID, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader.ID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader.ID, logLength, nil, gl.Str(log))

		return fmt.Errorf("shader compilation failed.\n** Source: %s**\n%v\n** Log: **\n%v", path, source, log)
	}

	return nil
}
