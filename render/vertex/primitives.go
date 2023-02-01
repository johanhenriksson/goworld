package vertex

import "github.com/vkngwrapper/core/v2/core1_0"

type Primitive core1_0.PrimitiveTopology

const (
	Triangles Primitive = Primitive(core1_0.PrimitiveTopologyTriangleList)
	Lines               = Primitive(core1_0.PrimitiveTopologyLineList)
	Points              = Primitive(core1_0.PrimitiveTopologyPointList)
)
