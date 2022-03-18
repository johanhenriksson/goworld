package pipeline

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/vertex"
	vk "github.com/vulkan-go/vulkan"
)

type Args struct {
	Pass     renderpass.T
	Subpass  string
	Layout   Layout
	Shader   Shader
	Pointers vertex.Pointers

	Primitive       vertex.Primitive
	PolygonFillMode vk.PolygonMode
	CullMode        vk.CullModeFlagBits

	DepthTest  bool
	DepthWrite bool
	DepthFunc  vk.CompareOp

	StencilTest bool
}

func (args *Args) defaults() {
	if args.DepthFunc == 0 {
		args.DepthFunc = vk.CompareOpLessOrEqual
	}
	if args.Primitive == 0 {
		args.Primitive = vertex.Triangles
	}
}

type PushConstant struct {
	Offset int
	Size   int
	Stages vk.ShaderStageFlagBits
}
