package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.DescriptorSetLayout]
}

type layout struct {
	ptr    vk.DescriptorSetLayout
	device device.T
}

type Binding struct {
	Binding int
	Type    vk.DescriptorType
	Count   int
	Stages  vk.ShaderStageFlags
}

func New(device device.T, bindings []Binding) T {
	info := vk.DescriptorSetLayoutCreateInfo{
		SType:        vk.StructureTypeDescriptorSetLayoutCreateInfo,
		BindingCount: uint32(len(bindings)),
		PBindings: util.Map(bindings, func(i int, layout Binding) vk.DescriptorSetLayoutBinding {
			return vk.DescriptorSetLayoutBinding{
				Binding:         uint32(layout.Binding),
				DescriptorType:  layout.Type,
				DescriptorCount: uint32(layout.Count),
				StageFlags:      layout.Stages,
			}
		}),
	}
	var ptr vk.DescriptorSetLayout
	vk.CreateDescriptorSetLayout(device.Ptr(), &info, nil, &ptr)

	return &layout{
		ptr:    ptr,
		device: device,
	}
}

func (l *layout) Ptr() vk.DescriptorSetLayout {
	return l.ptr
}

func (l *layout) Destroy() {
	vk.DestroyDescriptorSetLayout(l.device.Ptr(), l.ptr, nil)
	l.ptr = nil
}
