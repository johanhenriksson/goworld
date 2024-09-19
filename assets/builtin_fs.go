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

func (bfs *builtinFilesystem) Read(key string) ([]byte, error) {
	file, err := builtinFs.Open("builtin/" + key)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("builtin asset %s %w", key, fs.ErrNotFound)
		}
	}
	return io.ReadAll(file)
}

func (_ *builtinFilesystem) Write(key string, data []byte) error {
	return fmt.Errorf("%w: cant write to builtin file system", fs.ErrImmutable)
}
