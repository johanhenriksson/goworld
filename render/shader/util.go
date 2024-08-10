package shader

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"unsafe"

	"github.com/johanhenriksson/goworld/assets"
)

var ErrCompileFailed = errors.New("shader compilation error")

var includePattern = regexp.MustCompile(`(?m:^#include\s+\"[^\"]*\"\s*$)`)
var includeFilePattern = regexp.MustCompile(`#include\s+\"([^\"]*)\"`)

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
	source, err := assets.Read(spvPath)
	if errors.Is(err, assets.ErrNotFound) {
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

	source, err := LoadSource(path, []string{"shaders"})
	if err != nil {
		return nil, err
	}

	// todo: check for glslc
	bytecode := &bytes.Buffer{}
	errors := &bytes.Buffer{}
	args := []string{
		stageflag,
		"-O",      // optimize SPIR-V
		"-o", "-", // output file: standard out
		"-", // input file: standard in
	}
	cmd := exec.Command("glslc", args...)
	cmd.Stdin = bytes.NewBuffer(source)
	cmd.Stdout = bytecode
	cmd.Stderr = errors

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

func LoadSource(path string, includePaths []string) ([]byte, error) {
	source, err := assets.Read(path)
	if errors.Is(err, assets.ErrNotFound) {
		for _, includePath := range includePaths {
			includePathFile := filepath.Join(includePath, path)
			source, err = assets.Read(includePathFile)
			if errors.Is(err, assets.ErrNotFound) {
				continue
			} else if err != nil {
				return nil, err
			}
		}
	}
	if err != nil {
		return nil, err
	}

	// implement #include logic
	for {
		// find the next include statement
		include := includePattern.FindIndex(source)
		if include == nil {
			break
		}

		// extract the file name
		includeStatement := string(source[include[0]:include[1]])
		includeFile := includeFilePattern.FindStringSubmatch(includeStatement)
		if len(includeFile) != 2 {
			return nil, fmt.Errorf("invalid include statement: %s", includeStatement)
		}

		// recursively load the included file
		includeSource, err := LoadSource(includeFile[1], includePaths)
		if err != nil {
			return nil, err
		}

		// insert the included file into the source
		source = bytes.Replace(source, []byte(includeStatement), includeSource, 1)
	}

	// return preprocessed source
	return source, nil
}
