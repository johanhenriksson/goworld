package assets

import (
	"embed"
	"fmt"
	"io"
)

//go:embed builtin/*
var builtinFs embed.FS

var BuiltinFilesystem = &builtinFilesystem{}

type builtinFilesystem struct{}

func (fs *builtinFilesystem) Read(key string) ([]byte, error) {
	file, err := builtinFs.Open("builtin/" + key)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(file)
}

func (fs *builtinFilesystem) Write(key string, data []byte) error {
	return fmt.Errorf("cant write to immutable file system")
}
