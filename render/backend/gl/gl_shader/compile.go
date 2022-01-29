package gl_shader

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/shader"
)

// CompileShader compiles a set of GLSL files into a linked shader program.
// Filenames ending in vs, fs, gs indicate vertex, fragment and geometry shaders.
func CompileShader(name string, fileNames ...string) shader.T {
	shader := New(name)
	for _, fileName := range fileNames {
		stage := StageFromFile(fileName)
		shader.Attach(stage)
	}
	shader.Link()
	return shader
}

// CompileStageFromFile compiles and returns a shader from the given source file
// Panics on compilation errors
func CompileStageFromFile(kind shader.StageType, path string) shader.Stage {
	s := NewStage(kind)
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
