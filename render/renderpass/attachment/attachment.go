package attachment

import (
	vk "github.com/vulkan-go/vulkan"
)

type Name string

type T interface {
	Name() Name
	Allocator() Allocator
	Clear() vk.ClearValue
	Format() vk.Format
	Usage() vk.ImageUsageFlagBits
	Description() vk.AttachmentDescription
	Blend() Blend
}

type BlendOp struct {
	Operation vk.BlendOp
	SrcFactor vk.BlendFactor
	DstFactor vk.BlendFactor
}

type Blend struct {
	Enabled bool
	Color   BlendOp
	Alpha   BlendOp
}

type attachment struct {
	name   Name
	alloc  Allocator
	clear  vk.ClearValue
	desc   vk.AttachmentDescription
	blend  Blend
	format vk.Format
	usage  vk.ImageUsageFlagBits
}

func (a *attachment) Description() vk.AttachmentDescription {
	return a.desc
}

func (a *attachment) Name() Name                   { return a.name }
func (a *attachment) Allocator() Allocator         { return a.alloc }
func (a *attachment) Clear() vk.ClearValue         { return a.clear }
func (a *attachment) Blend() Blend                 { return a.blend }
func (a *attachment) Format() vk.Format            { return a.format }
func (a *attachment) Usage() vk.ImageUsageFlagBits { return a.usage }
