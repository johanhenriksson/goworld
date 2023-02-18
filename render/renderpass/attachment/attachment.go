package attachment

import (
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Name string

type T interface {
	Name() Name
	Image() Image
	Clear() core1_0.ClearValue
	Description() core1_0.AttachmentDescription
	Blend() Blend
}

type BlendOp struct {
	Operation core1_0.BlendOp
	SrcFactor core1_0.BlendFactor
	DstFactor core1_0.BlendFactor
}

type Blend struct {
	Enabled bool
	Color   BlendOp
	Alpha   BlendOp
}

type attachment struct {
	name  Name
	image Image
	clear core1_0.ClearValue
	desc  core1_0.AttachmentDescription
	blend Blend
}

func (a *attachment) Description() core1_0.AttachmentDescription {
	return a.desc
}

func (a *attachment) Name() Name                { return a.name }
func (a *attachment) Image() Image              { return a.image }
func (a *attachment) Clear() core1_0.ClearValue { return a.clear }
func (a *attachment) Blend() Blend              { return a.blend }
