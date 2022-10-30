package renderpass

import (
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	vk "github.com/vulkan-go/vulkan"
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

	Flags         vk.DependencyFlagBits
	SrcStageMask  vk.PipelineStageFlagBits
	SrcAccessMask vk.AccessFlagBits
	DstStageMask  vk.PipelineStageFlagBits
	DstAccessMask vk.AccessFlagBits
}
