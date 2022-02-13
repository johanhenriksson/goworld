package image

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Image]

	Memory() device.Memory
	View(format vk.Format, mask vk.ImageAspectFlags) View
}

type image struct {
	ptr    vk.Image
	device device.T
	memory device.Memory
}

func New2D(device device.T, width, height int, format vk.Format, usage vk.ImageUsageFlags) T {
	queueIdx := device.GetQueueFamilyIndex(vk.QueueFlags(vk.QueueGraphicsBit))
	info := vk.ImageCreateInfo{
		SType:     vk.StructureTypeImageCreateInfo,
		ImageType: vk.ImageType2d,
		Format:    format,
		Extent: vk.Extent3D{
			Width:  uint32(width),
			Height: uint32(height),
			Depth:  1,
		},
		MipLevels:   1,
		ArrayLayers: 1,
		Samples:     vk.SampleCountFlagBits(vk.SampleCount1Bit),
		Tiling:      vk.ImageTilingOptimal,
		// Usage: vk.ImageUsageFlags(
		// 	vk.ImageUsageDepthStencilAttachmentBit |
		// 		vk.ImageUsageTransferSrcBit),
		Usage:                 usage,
		SharingMode:           vk.SharingModeExclusive,
		QueueFamilyIndexCount: 1,
		PQueueFamilyIndices:   []uint32{uint32(queueIdx)},
		InitialLayout:         vk.ImageLayoutUndefined,
	}

	var ptr vk.Image
	vk.CreateImage(device.Ptr(), &info, nil, &ptr)

	var memreq vk.MemoryRequirements
	vk.GetImageMemoryRequirements(device.Ptr(), ptr, &memreq)
	memreq.Deref()

	mem := device.Allocate(memreq, vk.MemoryPropertyFlags(vk.MemoryPropertyDeviceLocalBit))
	vk.BindImageMemory(device.Ptr(), ptr, mem.Ptr(), vk.DeviceSize(0))

	return &image{
		ptr:    ptr,
		device: device,
		memory: mem,
	}
}

func Wrap(device device.T, ptr vk.Image) T {
	return &image{
		ptr:    ptr,
		device: device,
		memory: nil,
	}
}

func (i *image) Ptr() vk.Image {
	return i.ptr
}

func (i *image) Memory() device.Memory {
	return i.memory
}

func (i *image) Destroy() {
	if i.memory != nil {
		i.memory.Destroy()
	}

	vk.DestroyImage(i.device.Ptr(), i.ptr, nil)
	i.ptr = nil
}

func (i *image) View(format vk.Format, mask vk.ImageAspectFlags) View {
	info := vk.ImageViewCreateInfo{
		SType:    vk.StructureTypeImageViewCreateInfo,
		Image:    i.ptr,
		ViewType: vk.ImageViewType2d,
		Format:   format,
		SubresourceRange: vk.ImageSubresourceRange{
			AspectMask:     mask,
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	var ptr vk.ImageView
	vk.CreateImageView(i.device.Ptr(), &info, nil, &ptr)

	return &imgview{
		ptr:    ptr,
		device: i.device,
		image:  i,
		format: format,
	}
}
