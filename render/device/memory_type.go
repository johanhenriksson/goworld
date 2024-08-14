package device

import (
	"github.com/vkngwrapper/core/v2/core1_0"
)

type MemoryType int

const (
	MemoryTypeGPU     MemoryType = 1
	MemoryTypeShared  MemoryType = 2
	MemoryTypeCPU     MemoryType = 3
	MemoryTypeTexture MemoryType = 4
)

type memtype struct {
	Index int
	Flags core1_0.MemoryPropertyFlags
}

func getImageMemoryTypeBits(dev core1_0.Device) (uint32, error) {
	img, _, err := dev.CreateImage(nil, core1_0.ImageCreateInfo{
		ImageType:   core1_0.ImageType2D,
		Format:      core1_0.FormatR8G8B8A8UnsignedNormalized,
		Extent:      core1_0.Extent3D{Width: 1, Height: 1, Depth: 1},
		MipLevels:   1,
		ArrayLayers: 1,
		Samples:     1,
		Usage:       core1_0.ImageUsageColorAttachment | core1_0.ImageUsageSampled,
	})
	if err != nil {
		return 0, err
	}
	memReqs := img.MemoryRequirements()
	img.Destroy(nil)
	return memReqs.MemoryTypeBits, nil
}

func getBufferMemoryTypeBits(dev core1_0.Device) (uint32, error) {
	buf, _, err := dev.CreateBuffer(nil, core1_0.BufferCreateInfo{
		Size:        1,
		Usage:       core1_0.BufferUsageTransferSrc,
		SharingMode: core1_0.SharingModeExclusive,
	})
	memReqs := buf.MemoryRequirements()
	if err != nil {
		return 0, err
	}
	buf.Destroy(nil)
	return memReqs.MemoryTypeBits, nil
}

func pickPreferredMemoryType(
	types []core1_0.MemoryType,
	typeMask uint32,
	requiredFlags core1_0.MemoryPropertyFlags,
	preferredFlags core1_0.MemoryPropertyFlags,
) memtype {
	findMatchingType := func(flags core1_0.MemoryPropertyFlags) memtype {
		for i, kind := range types {
			typeBits := uint32(1 << i)
			if typeMask&typeBits != typeBits {
				continue
			}

			if kind.PropertyFlags&flags != flags {
				continue
			}

			return memtype{i, kind.PropertyFlags}
		}
		return memtype{-1, 0}
	}

	if t := findMatchingType(requiredFlags | preferredFlags); t.Index != -1 {
		return t
	}

	return findMatchingType(requiredFlags)
}
