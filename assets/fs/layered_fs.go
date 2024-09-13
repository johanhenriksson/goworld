package fs

import (
	"fmt"
	"slices"
)

type Layered struct {
	layers []Filesystem
}

var _ Filesystem = (*Layered)(nil)

func NewLayered(layers ...Filesystem) *Layered {
	return &Layered{layers: layers}
}

// Push adds a layer to the top of the filesystem stack
func (fs *Layered) Push(layer Filesystem) {
	fs.layers = slices.Insert(fs.layers, 0, layer)
}

// Pop removes the top layer from the filesystem stack
func (fs *Layered) Pop() {
	fs.layers = slices.Delete(fs.layers, 0, 1)
}

func (fs *Layered) Read(key string) ([]byte, error) {
	for _, layer := range fs.layers {
		data, err := layer.Read(key)
		if err == nil {
			return data, nil
		}
	}
	return nil, fmt.Errorf("asset %s %w", key, ErrNotFound)
}

func (fs *Layered) Write(key string, data []byte) error {
	if len(fs.layers) == 0 {
		return fmt.Errorf("no layers in filesystem")
	}
	// write to the top layer
	return fs.layers[0].Write(key, data)
}
