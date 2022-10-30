package attachment

import (
	"github.com/johanhenriksson/goworld/render/image"
	vk "github.com/vulkan-go/vulkan"
)

type Name string

type T interface {
	Name() Name
	Image(int) image.T
	View(int) image.View
	Clear() vk.ClearValue
	Description() vk.AttachmentDescription
	Destroy()
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
	name     Name
	view     []image.View
	image    []image.T
	clear    vk.ClearValue
	desc     vk.AttachmentDescription
	imgowner bool
	blend    Blend
}

func (a *attachment) Description() vk.AttachmentDescription {
	return a.desc
}

func (a *attachment) Name() Name                { return a.name }
func (a *attachment) Image(frame int) image.T   { return a.image[frame%len(a.image)] }
func (a *attachment) View(frame int) image.View { return a.view[frame%len(a.view)] }
func (a *attachment) Clear() vk.ClearValue      { return a.clear }
func (a *attachment) Blend() Blend              { return a.blend }

func (a *attachment) Destroy() {
	for i := range a.image {
		a.view[i].Destroy()
		if a.imgowner {
			a.image[i].Destroy()
		}
	}
}
