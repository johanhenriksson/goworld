package pipeline

import (
	"io/ioutil"
	"os"
	"unsafe"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Shader interface {
	device.Resource[vk.ShaderModule]

	Entrypoint() string
	Stage() vk.ShaderStageFlagBits
}

type shader struct {
	device device.T
	ptr    vk.ShaderModule
	stage  vk.ShaderStageFlagBits
}

func NewShader(device device.T, path string, stage vk.ShaderStageFlagBits) Shader {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	info := vk.ShaderModuleCreateInfo{
		SType:    vk.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint(len(data)),
		PCode:    sliceUint32(data),
	}

	var ptr vk.ShaderModule
	vk.CreateShaderModule(device.Ptr(), &info, nil, &ptr)

	return &shader{
		device: device,
		ptr:    ptr,
		stage:  stage,
	}
}

func (s *shader) Ptr() vk.ShaderModule {
	return s.ptr
}

func (s *shader) Stage() vk.ShaderStageFlagBits {
	return s.stage
}

func (s *shader) Entrypoint() string {
	return "main"
}

func (s *shader) Destroy() {
	vk.DestroyShaderModule(s.device.Ptr(), s.ptr, nil)
	s.ptr = nil
}

// Disgusting hack that reinterprets a byte slice as a slice of uint32
func sliceUint32(data []byte) []uint32 {
	type sliceHeader struct {
		Data uintptr
		Len  int
		Cap  int
	}
	const m = 0x7fffffff
	return (*[m / 4]uint32)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&data)).Data))[:len(data)/4]
}
