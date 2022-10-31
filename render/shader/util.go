package shader

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"unsafe"
)

var ErrCompileFailed = errors.New("compilation failed")

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
