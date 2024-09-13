package assets

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/johanhenriksson/goworld/assets/fs"
)

//go:embed builtin/*
var builtinFs embed.FS

var BuiltinFilesystem fs.Filesystem = &builtinFilesystem{}

type builtinFilesystem struct{}

func (fs *builtinFilesystem) Read(key string) ([]byte, error) {
	file, err := builtinFs.Open("builtin/" + key)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("asset %s %w", key, ErrNotFound)
		}
	}
	return io.ReadAll(file)
}

func (fs *builtinFilesystem) Write(key string, data []byte) error {
	return fmt.Errorf("cant write to immutable file system")
}
