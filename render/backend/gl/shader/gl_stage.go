package shader

import (
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

// NewStage creates a new empty shader
func NewStage(shaderType shader.StageType) shader.Stage {
	id := gl.CreateShader(uint32(shaderType))
	return &stage{
		id:       id,
		stype:    shaderType,
		compiled: false,
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
