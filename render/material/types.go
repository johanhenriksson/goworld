package material

import "github.com/johanhenriksson/goworld/render/vertex"

// todo: this is rather implementation specific and likely
// does not belong in the render package

func StandardDeferred() *Def {
	return &Def{
		Pass:         Deferred,
		Shader:       "deferred/textured",
		VertexFormat: vertex.T{},
		DepthTest:    true,
		DepthWrite:   true,
		DepthClamp:   false,
		Primitive:    vertex.Triangles,
		CullMode:     vertex.CullBack,
	}
}

func StandardForward() *Def {
	return &Def{
		Pass:         Forward,
		Shader:       "forward/textured",
		VertexFormat: vertex.T{},
		DepthTest:    true,
		DepthWrite:   true,
		DepthClamp:   false,
		Primitive:    vertex.Triangles,
		CullMode:     vertex.CullBack,
	}
}

func ColoredForward() *Def {
	return &Def{
		Pass:         Forward,
		Shader:       "forward/color",
		VertexFormat: vertex.C{},
		DepthTest:    true,
		DepthWrite:   true,
		DepthClamp:   false,
		Primitive:    vertex.Triangles,
		CullMode:     vertex.CullBack,
	}
}

func Lines() *Def {
	return &Def{
		Shader:       "lines",
		Pass:         "lines",
		VertexFormat: vertex.C{},
		Primitive:    vertex.Lines,
		DepthTest:    true,
		DepthWrite:   false,
		DepthClamp:   false,
		CullMode:     vertex.CullNone,
	}
}
