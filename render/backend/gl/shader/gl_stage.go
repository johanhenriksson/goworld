package shader

import (
	"fmt"
	"io/ioutil"

	"github.com/johanhenriksson/goworld/render/backend/gl"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/util"
)

// stage represents a shader part of a GLSL program
type stage struct {
	id       shader.StageID
	stype    shader.StageType
	compiled bool
}

// CreateStage creates a new empty shader
func CreateStage(shaderType shader.StageType) shader.Stage {
	id := gl.CreateShader(uint32(shaderType))
	return &stage{
		id:       id,
		stype:    shaderType,
		compiled: false,
	}
}

// VertexShader compiles and returns a vertex shader from the given source file
// Panics on compilation errors
func CompileStageFromFile(kind shader.StageType, path string) shader.Stage {
	s := CreateStage(kind)
	err := s.CompileFile(path)
	if err != nil {
		panic(err)
	}
	return s
}

func StageFromFile(fileName string) shader.Stage {
	if len(fileName) < 3 {
		panic(fmt.Errorf("invalid shader filename: %s", fileName))
	}
	kind := fileName[len(fileName)-3:]
	switch kind {
	case ".fs":
		return CompileStageFromFile(shader.FragmentShader, fmt.Sprintf("%s.glsl", fileName))
	case ".vs":
		return CompileStageFromFile(shader.VertexShader, fmt.Sprintf("%s.glsl", fileName))
	case ".gs":
		return CompileStageFromFile(shader.GeometryShader, fmt.Sprintf("%s.glsl", fileName))
	default:
		panic(fmt.Errorf("invalid shader type %s: %s", kind, fileName))
	}
}

func (stage *stage) ID() shader.StageID {
	return stage.id
}

// CompileFile loads and compiles source code from the given file path
func (stage *stage) CompileFile(path string) error {
	source, err := ioutil.ReadFile(util.ExePath + path)
	if err != nil {
		return err
	}
	return stage.Compile(string(source), path)
}

// Compile a shader from a source string
func (stage *stage) Compile(source, path string) error {
	return gl.CompileShader(stage.id, source, path)
}
