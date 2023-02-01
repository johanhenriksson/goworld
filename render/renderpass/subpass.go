package renderpass

import (
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Name string

const ExternalSubpass Name = "external"

type Subpass struct {
	index int

	Name             Name
	Depth            bool
	ColorAttachments []attachment.Name
	InputAttachments []attachment.Name
}

func (s *Subpass) Index() int {
	return s.index
}

type SubpassDependency struct {
	Src Name
	Dst Name

	Flags         core1_0.DependencyFlags
	SrcStageMask  core1_0.PipelineStageFlags
	SrcAccessMask core1_0.AccessFlags
	DstStageMask  core1_0.PipelineStageFlags
	DstAccessMask core1_0.AccessFlags
}
