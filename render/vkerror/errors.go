package vkerror

import (
	"errors"
	"fmt"

	vk "github.com/vulkan-go/vulkan"
)

var ErrOutOfHostMemory = errors.New("out of host memory")
var ErrOutOfDeviceMemory = errors.New("out of device memory")

func FromResult(result vk.Result) error {
	switch result {
	case vk.Success:
		return nil

	case vk.ErrorOutOfHostMemory:
		return ErrOutOfHostMemory

	case vk.ErrorOutOfDeviceMemory:
		return ErrOutOfDeviceMemory

	default:
		return fmt.Errorf("unmapped Vulkan error: %d", result)
	}
}
