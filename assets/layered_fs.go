package assets

import (
	"fmt"
	"slices"
)

type LayeredFilesystem struct {
	layers []Server
}

func NewLayeredFilesystem(layers ...Server) *LayeredFilesystem {
	return &LayeredFilesystem{layers: layers}
}

// Push adds a layer to the top of the filesystem stack
func (fs *LayeredFilesystem) Push(layer Server) {
	fs.layers = slices.Insert(fs.layers, 0, layer)
}

// Pop removes the top layer from the filesystem stack
func (fs *LayeredFilesystem) Pop() {
	fs.layers = slices.Delete(fs.layers, 0, 1)
}

func (fs *LayeredFilesystem) Read(key string) ([]byte, error) {
	for _, layer := range fs.layers {
		data, err := layer.Read(key)
		if err == nil {
			return data, nil
		}
	}
	return nil, ErrNotFound
}

func (fs *LayeredFilesystem) Write(key string, data []byte) error {
	if len(fs.layers) == 0 {
		return fmt.Errorf("no layers in filesystem")
	}
	// write to the top layer
	return fs.layers[0].Write(key, data)
}
