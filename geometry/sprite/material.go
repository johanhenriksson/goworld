package sprite

import (
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

func Material() *material.Def {
	return &material.Def{
		Pass:         material.Forward,
		Shader:       "forward/sprite",
		VertexFormat: vertex.Vertex{},
		DepthTest:    true,
		DepthWrite:   true,
		DepthFunc:    core1_0.CompareOpLessOrEqual,
		Primitive:    vertex.Triangles,
		CullMode:     vertex.CullNone,
		Transparent:  true,
	}
}
