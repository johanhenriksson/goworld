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

func CompileShaderStage(fileName string) *ShaderStage {
	if len(fileName) < 3 {
		panic(fmt.Errorf("invalid shader filename: %s", fileName))
	}
	kind := fileName[len(fileName)-3:]
	switch kind {
	case ".fs":
		return FragmentShader(fmt.Sprintf("%s.glsl", fileName))
	case ".vs":
		return VertexShader(fmt.Sprintf("%s.glsl", fileName))
	case ".gs":
		return GeometryShader(fmt.Sprintf("%s.glsl", fileName))
	default:
		panic(fmt.Errorf("invalid shader type %s: %s", kind, fileName))
	}
}

// CompileFile loads and compiles source code from the given file path
func (stage *ShaderStage) CompileFile(path string) error {
	source, err := ioutil.ReadFile(util.ExePath + path)
	if err != nil {
		return err
	}
	return stage.Compile(string(source), path)
}

// Compile a shader from a source string
func (stage *ShaderStage) Compile(source, path string) error {
	csource, free := util.GLString(source)
	gl.ShaderSource(stage.ID, 1, csource, nil)
	gl.CompileShader(stage.ID)
	free()

	/* Check compilation status */
	var status int32
	gl.GetShaderiv(stage.ID, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(stage.ID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(stage.ID, logLength, nil, gl.Str(log))

		return fmt.Errorf("shader compilation failed.\n** Source: %s**\n%v\n** Log: **\n%v", path, source, log)
	}

	return nil
}
