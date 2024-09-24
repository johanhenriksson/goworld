package material

import (
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// todo: this is rather implementation specific and likely
// does not belong in the render package

func StandardDeferred() *Def {
	return &Def{
		Pass:         Deferred,
		Shader:       "deferred/textured",
		VertexFormat: vertex.Vertex{},
		DepthTest:    true,
		DepthWrite:   true,
		Primitive:    vertex.Triangles,
		CullMode:     vertex.CullBack,
	}
}

func StandardForward() *Def {
	return &Def{
		Pass:         Forward,
		Shader:       "forward/textured",
		VertexFormat: vertex.Vertex{},
		DepthTest:    true,
		DepthWrite:   true,
		DepthFunc:    core1_0.CompareOpLessOrEqual,
		Primitive:    vertex.Triangles,
		CullMode:     vertex.CullBack,
		Transparent:  false,
	}
}

func TransparentForward() *Def {
	return &Def{
		Pass:         Forward,
		Shader:       "forward/textured",
		VertexFormat: vertex.Vertex{},
		DepthTest:    true,
		DepthWrite:   true,
		DepthFunc:    core1_0.CompareOpLessOrEqual,
		Primitive:    vertex.Triangles,
		CullMode:     vertex.CullBack,
		Transparent:  true,
	}
}

func ColoredForward() *Def {
	return &Def{
		Pass:         Forward,
		Shader:       "forward/color",
		VertexFormat: vertex.Vertex{},
		DepthTest:    true,
		DepthWrite:   true,
		DepthFunc:    core1_0.CompareOpLessOrEqual,
		Primitive:    vertex.Triangles,
		CullMode:     vertex.CullBack,
	}
}

func Lines() *Def {
	return &Def{
		Shader:       "lines",
		Pass:         "lines",
		VertexFormat: vertex.Vertex{},
		Primitive:    vertex.Lines,
		DepthTest:    true,
		DepthWrite:   false,
		DepthFunc:    core1_0.CompareOpLessOrEqual,
		CullMode:     vertex.CullNone,
	}
}
