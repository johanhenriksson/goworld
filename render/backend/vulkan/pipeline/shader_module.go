package pipeline

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"unsafe"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type ShaderModule interface {
	device.Resource[vk.ShaderModule]

	Entrypoint() string
	Stage() vk.ShaderStageFlagBits
}

type shader_module struct {
	device device.T
	ptr    vk.ShaderModule
	stage  vk.ShaderStageFlagBits
}

func NewShaderModule(device device.T, path string, stage vk.ShaderStageFlagBits) ShaderModule {
	bytecode, err := LoadOrCompile(path)
	if err != nil {
		panic(err)
	}

	info := vk.ShaderModuleCreateInfo{
		SType:    vk.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint(len(bytecode)),
		PCode:    sliceUint32(bytecode),
	}

	var ptr vk.ShaderModule
	vk.CreateShaderModule(device.Ptr(), &info, nil, &ptr)

	return &shader_module{
		device: device,
		ptr:    ptr,
		stage:  stage,
	}
}

func (s *shader_module) Ptr() vk.ShaderModule {
	return s.ptr
}

func (s *shader_module) Stage() vk.ShaderStageFlagBits {
	return s.stage
}

func (s *shader_module) Entrypoint() string {
	return "main"
}

func (s *shader_module) Destroy() {
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

var ErrCompileFailed = errors.New("compilation failed")

func LoadOrCompile(path string) ([]byte, error) {
	spvPath := fmt.Sprintf("%s.spv", path)
	fp, err := os.Open(spvPath)
	if errors.Is(err, os.ErrNotExist) {
		return Compile(path)
	}
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	log.Println("loading shader", path)
	return ioutil.ReadAll(fp)
}

func Compile(path string) ([]byte, error) {
	// check for glslc
	bytecode := &bytes.Buffer{}
	errors := &bytes.Buffer{}
	cmd := exec.Command("glslc", path, "-o", "-", "-O")
	cmd.Stdout = bytecode
	cmd.Stderr = errors

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w in shader %s:\n%s",
			ErrCompileFailed,
			path,
			errors.String())
	}

	log.Println("shader compiled successfully:", path)
	return bytecode.Bytes(), nil
}
