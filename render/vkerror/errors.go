package vkerror

import (
	"errors"
	"fmt"

	"github.com/vkngwrapper/core/v2/common"
	"github.com/vkngwrapper/core/v2/core1_0"
)

var ErrOutOfHostMemory = errors.New("out of host memory")
var ErrOutOfDeviceMemory = errors.New("out of device memory")

func FromResult(result common.VkResult) error {
	switch result {
	case core1_0.VKSuccess:
		return nil

	case core1_0.VKErrorOutOfHostMemory:
		return ErrOutOfHostMemory

	case core1_0.VKErrorOutOfDeviceMemory:
		return ErrOutOfDeviceMemory

	default:
		return fmt.Errorf("unmapped Vulkan error: %d", result)
	}
}
