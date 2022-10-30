package renderpass

import (
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
)

type Args struct {
	Frames int
	Width  int
	Height int

	ColorAttachments []attachment.Color
	DepthAttachment  *attachment.Depth

	Subpasses    []Subpass
	Dependencies []SubpassDependency
}
