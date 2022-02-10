package vulkan

type T interface {
	Resize(int, int)
}

type vulkan struct {
}

func (v *vulkan) Resize(width, height int) {
	// recreate swapchain?

}
