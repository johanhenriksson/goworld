package image

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vkerror"

	vk "github.com/vulkan-go/vulkan"
)

// Represents vk.NullImage
var Nil T = &image{
	ptr:    vk.NullImage,
	device: device.Nil,
	memory: device.NilMemory,
}

type T interface {
	device.Resource[vk.Image]

	Memory() device.Memory
	View(format vk.Format, mask vk.ImageAspectFlags) (View, error)
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

func New2D(device device.T, width, height int, format vk.Format, usage vk.ImageUsageFlags) (T, error) {
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

func New(device device.T, args Args) (T, error) {
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
	result := vk.CreateImage(device.Ptr(), &info, nil, &ptr)
	if result != vk.Success {
		return nil, vkerror.FromResult(result)
	}

	var memreq vk.MemoryRequirements
	vk.GetImageMemoryRequirements(device.Ptr(), ptr, &memreq)
	memreq.Deref()

	mem := device.Allocate(memreq, vk.MemoryPropertyFlags(args.Memory))
	result = vk.BindImageMemory(device.Ptr(), ptr, mem.Ptr(), vk.DeviceSize(0))
	if result != vk.Success {
		// clean up
		vk.DestroyImage(device.Ptr(), ptr, nil)
		return nil, vkerror.FromResult(result)
	}

	return &image{
		Args:   args,
		ptr:    ptr,
		device: device,
		memory: mem,
	}, nil
}

func Wrap(dev device.T, ptr vk.Image) T {
	return &image{
		ptr:    ptr,
		device: dev,
		memory: device.NilMemory,
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
	if i.memory != device.NilMemory {
		i.memory.Destroy()
		if i.ptr != vk.NullImage {
			vk.DestroyImage(i.device.Ptr(), i.ptr, nil)
		}
	}
	i.ptr = vk.NullImage
	i.memory = device.NilMemory
	i.device = device.Nil
}

func (i *image) View(format vk.Format, mask vk.ImageAspectFlags) (View, error) {
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
	result := vk.CreateImageView(i.device.Ptr(), &info, nil, &ptr)
	if result != vk.Success {
		return nil, vkerror.FromResult(result)
	}

	return &imgview{
		ptr:    ptr,
		device: i.device,
		image:  i,
		format: format,
	}, nil
}
