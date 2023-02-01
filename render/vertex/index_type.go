package vertex

import "github.com/vkngwrapper/core/v2/core1_0"

func IndexType(size int) core1_0.IndexType {
	switch size {
	case 2:
		return core1_0.IndexTypeUInt16
	case 4:
		return core1_0.IndexTypeUInt32
	default:
		panic("illegal index size")
	}
}
