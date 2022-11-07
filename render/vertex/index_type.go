package vertex

import (
	vk "github.com/vulkan-go/vulkan"
)

func IndexType(size int) vk.IndexType {
	switch size {
	case 2:
		return vk.IndexTypeUint16
	case 4:
		return vk.IndexTypeUint32
	default:
		panic("illegal index size")
	}
}
