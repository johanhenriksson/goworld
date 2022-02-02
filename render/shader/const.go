package shader

import (
	ogl "github.com/go-gl/gl/v4.1-core/gl"
)

// StageType indicates the type of shader program
type StageType uint32

// VertexShader is a Vertex Shader
const VertexShader StageType = ogl.VERTEX_SHADER

// FragmentShader is a Fragment Shader
const FragmentShader StageType = ogl.FRAGMENT_SHADER

const GeometryShader StageType = ogl.GEOMETRY_SHADER
