package renderpass

import (
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
)

type Args struct {
	ColorAttachments []attachment.Color
	DepthAttachment  *attachment.Depth

	Subpasses    []Subpass
	Dependencies []SubpassDependency
}
