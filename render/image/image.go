package image

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vkerror"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type T interface {
	device.Resource[core1_0.Image]

	Memory() device.Memory
	View(format core1_0.Format, mask core1_0.ImageAspectFlags) (View, error)
	Width() int
	Height() int
	Format() core1_0.Format
	Size() vec3.T
}

type image struct {
	Args
	ptr    core1_0.Image
	device device.T
	memory device.Memory
}

type Args struct {
	Type    core1_0.ImageType
	Width   int
	Height  int
	Depth   int
	Layers  int
	Levels  int
	Format  core1_0.Format
	Usage   core1_0.ImageUsageFlags
	Tiling  core1_0.ImageTiling
	Sharing core1_0.SharingMode
	Layout  core1_0.ImageLayout
	Memory  core1_0.MemoryPropertyFlags
}

func New2D(device device.T, width, height int, format core1_0.Format, usage core1_0.ImageUsageFlags) (T, error) {
	return New(device, Args{
		Type:    core1_0.ImageType2D,
		Width:   width,
		Height:  height,
		Depth:   1,
		Layers:  1,
		Levels:  1,
		Format:  format,
		Usage:   usage,
		Tiling:  core1_0.ImageTilingOptimal,
		Sharing: core1_0.SharingModeExclusive,
		Layout:  core1_0.ImageLayoutUndefined,
		Memory:  core1_0.MemoryPropertyDeviceLocal,
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

	queueIdx := device.GetQueueFamilyIndex(core1_0.QueueGraphics)
	info := core1_0.ImageCreateInfo{
		ImageType: args.Type,
		Format:    args.Format,
		Extent: core1_0.Extent3D{
			Width:  args.Width,
			Height: args.Height,
			Depth:  args.Depth,
		},
		MipLevels:          args.Levels,
		ArrayLayers:        args.Layers,
		Samples:            core1_0.Samples1,
		Tiling:             args.Tiling,
		Usage:              core1_0.ImageUsageFlags(args.Usage),
		SharingMode:        args.Sharing,
		QueueFamilyIndices: []uint32{uint32(queueIdx)},
		InitialLayout:      args.Layout,
	}

	ptr, result, err := device.Ptr().CreateImage(nil, info)
	if err != nil {
		return nil, err
	}
	if result != core1_0.VKSuccess {
		return nil, vkerror.FromResult(result)
	}

	memreq := ptr.MemoryRequirements()

	mem := device.Allocate(core1_0.MemoryRequirements{
		Size:           int(memreq.Size),
		Alignment:      int(memreq.Alignment),
		MemoryTypeBits: memreq.MemoryTypeBits,
	}, core1_0.MemoryPropertyFlags(args.Memory))
	result, err = ptr.BindImageMemory(mem.Ptr(), 0)
	if err != nil {
		ptr.Destroy(nil)
		mem.Destroy()
		return nil, err
	}
	if result != core1_0.VKSuccess {
		ptr.Destroy(nil)
		mem.Destroy()
		return nil, vkerror.FromResult(result)
	}

	return &image{
		Args:   args,
		ptr:    ptr,
		device: device,
		memory: mem,
	}, nil
}

func Wrap(dev device.T, ptr core1_0.Image, args Args) T {
	return &image{
		ptr:    ptr,
		device: dev,
		memory: nil,
		Args:   args,
	}
}

func (i *image) Ptr() core1_0.Image {
	return i.ptr
}

func (i *image) Memory() device.Memory {
	return i.memory
}

func (i *image) Width() int             { return i.Args.Width }
func (i *image) Height() int            { return i.Args.Height }
func (i *image) Format() core1_0.Format { return i.Args.Format }

func (i *image) Size() vec3.T {
	return vec3.T{
		X: float32(i.Args.Width),
		Y: float32(i.Args.Height),
		Z: float32(i.Args.Depth),
	}
}

func (i *image) Destroy() {
	if i.memory != nil {
		i.memory.Destroy()
		if i.ptr != nil {
			i.ptr.Destroy(nil)
		}
	}
	i.ptr = nil
	i.memory = nil
	i.device = nil
}

func (i *image) View(format core1_0.Format, mask core1_0.ImageAspectFlags) (View, error) {
	info := core1_0.ImageViewCreateInfo{
		Image:    i.ptr,
		ViewType: core1_0.ImageViewType2D,
		Format:   format,
		SubresourceRange: core1_0.ImageSubresourceRange{
			AspectMask:     mask,
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	ptr, result, err := i.device.Ptr().CreateImageView(nil, info)
	if err != nil {
		return nil, err
	}
	if result != core1_0.VKSuccess {
		return nil, vkerror.FromResult(result)
	}

	return &imgview{
		ptr:    ptr,
		device: i.device,
		image:  i,
		format: format,
	}, nil
}
