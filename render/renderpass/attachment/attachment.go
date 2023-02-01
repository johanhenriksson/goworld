package attachment

import (
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Name string

type T interface {
	Name() Name
	Allocator() Allocator
	Clear() core1_0.ClearValue
	Format() core1_0.Format
	Usage() core1_0.ImageUsageFlags
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
	name   Name
	alloc  Allocator
	clear  core1_0.ClearValue
	desc   core1_0.AttachmentDescription
	blend  Blend
	format core1_0.Format
	usage  core1_0.ImageUsageFlags
}

func (a *attachment) Description() core1_0.AttachmentDescription {
	return a.desc
}

func (a *attachment) Name() Name                     { return a.name }
func (a *attachment) Allocator() Allocator           { return a.alloc }
func (a *attachment) Clear() core1_0.ClearValue      { return a.clear }
func (a *attachment) Blend() Blend                   { return a.blend }
func (a *attachment) Format() core1_0.Format         { return a.format }
func (a *attachment) Usage() core1_0.ImageUsageFlags { return a.usage }
