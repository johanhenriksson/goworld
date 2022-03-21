package image

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Image]

	Memory() device.Memory
	View(format vk.Format, mask vk.ImageAspectFlags) View
	Width() int
	Height() int
	Format() vk.Format
}

type image struct {
	Args
	ptr    vk.Image
	device device.T
	memory device.Memory
}

type Args struct {
	Type    vk.ImageType
	Width   int
	Height  int
	Depth   int
	Layers  int
	Levels  int
	Format  vk.Format
	Usage   vk.ImageUsageFlagBits
	Tiling  vk.ImageTiling
	Sharing vk.SharingMode
	Layout  vk.ImageLayout
	Memory  vk.MemoryPropertyFlagBits
}

func New2D(device device.T, width, height int, format vk.Format, usage vk.ImageUsageFlags) T {
	return New(device, Args{
		Type:    vk.ImageType2d,
		Width:   width,
		Height:  height,
		Depth:   1,
		Layers:  1,
		Levels:  1,
		Format:  format,
		Usage:   vk.ImageUsageFlagBits(usage),
		Tiling:  vk.ImageTilingOptimal,
		Sharing: vk.SharingModeExclusive,
		Layout:  vk.ImageLayoutUndefined,
		Memory:  vk.MemoryPropertyDeviceLocalBit,
	})
}

func New(device device.T, args Args) T {
	if args.Depth < 1 {
		args.Depth = 1
	}
	if args.Levels < 1 {
		args.Levels = 1
	}
	if args.Layers < 1 {
		args.Layers = 1
	}

	queueIdx := device.GetQueueFamilyIndex(vk.QueueFlags(vk.QueueGraphicsBit))
	info := vk.ImageCreateInfo{
		SType:     vk.StructureTypeImageCreateInfo,
		ImageType: args.Type,
		Format:    args.Format,
		Extent: vk.Extent3D{
			Width:  uint32(args.Width),
			Height: uint32(args.Height),
			Depth:  uint32(args.Depth),
		},
		MipLevels:             uint32(args.Levels),
		ArrayLayers:           uint32(args.Layers),
		Samples:               vk.SampleCountFlagBits(vk.SampleCount1Bit),
		Tiling:                args.Tiling,
		Usage:                 vk.ImageUsageFlags(args.Usage),
		SharingMode:           args.Sharing,
		QueueFamilyIndexCount: 1,
		PQueueFamilyIndices:   []uint32{uint32(queueIdx)},
		InitialLayout:         args.Layout,
	}

	var ptr vk.Image
	vk.CreateImage(device.Ptr(), &info, nil, &ptr)

	var memreq vk.MemoryRequirements
	vk.GetImageMemoryRequirements(device.Ptr(), ptr, &memreq)
	memreq.Deref()

	mem := device.Allocate(memreq, vk.MemoryPropertyFlags(args.Memory))
	vk.BindImageMemory(device.Ptr(), ptr, mem.Ptr(), vk.DeviceSize(0))

	return &image{
		Args:   args,
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

func (i *image) Width() int        { return i.Args.Width }
func (i *image) Height() int       { return i.Args.Height }
func (i *image) Format() vk.Format { return i.Args.Format }

func (i *image) Destroy() {
	if i.memory != nil {
		i.memory.Destroy()
	}

	if i.ptr != nil {
		vk.DestroyImage(i.device.Ptr(), i.ptr, nil)
		i.ptr = nil
	}
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
