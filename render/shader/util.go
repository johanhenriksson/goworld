package shader

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"unsafe"

	"github.com/johanhenriksson/goworld/assets"
)

var ErrCompileFailed = errors.New("shader compilation error")

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

func LoadOrCompile(path string, stage ShaderStage) ([]byte, error) {
	spvPath := fmt.Sprintf("%s.spv", path)
	source, err := assets.ReadAll(spvPath)
	if errors.Is(err, os.ErrNotExist) {
		return Compile(path, stage)
	}
	if err != nil {
		return nil, err
	}
	log.Println("loading shader", path)
	return source, nil
}

func Compile(path string, stage ShaderStage) ([]byte, error) {
	stageflag := ""
	switch stage {
	case StageFragment:
		stageflag = "-fshader-stage=fragment"
	case StageVertex:
		stageflag = "-fshader-stage=vertex"
	case StageCompute:
		stageflag = "-fshader-stage=compute"
	}

	source, err := assets.ReadAll(path)
	if err != nil {
		return nil, err
	}

	// todo: check for glslc
	includePath := filepath.Join(assets.Path, "shaders")
	bytecode := &bytes.Buffer{}
	errors := &bytes.Buffer{}
	args := []string{
		stageflag,
		"-O",              // optimize SPIR-V
		"-I", includePath, // include path
		"-o", "-", // output file: standard out
		"-", // input file: standard in
	}
	cmd := exec.Command("glslc", args...)
	cmd.Stdin = bytes.NewBuffer(source)
	cmd.Stdout = bytecode
	cmd.Stderr = errors
	cmd.Dir = assets.Path

	if err := cmd.Run(); err != nil {
		if errors.Len() > 0 {
			return nil, fmt.Errorf("%w in %s:\n%s",
				ErrCompileFailed,
				path,
				errors.String())
		}
		return nil, fmt.Errorf("%s in %s: %w", ErrCompileFailed, path, err)
	}

	log.Println("shader compiled successfully:", path)
	return bytecode.Bytes(), nil
}
