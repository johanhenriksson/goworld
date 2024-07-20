package assets

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Server interface {
	Read(key string) ([]byte, error)
	Write(key string, data []byte) error
}

type LocalFilesystem struct {
	root string
}

var _ Server = (*LocalFilesystem)(nil)

func NewLocalFilesystem(root string) *LocalFilesystem {
	return &LocalFilesystem{root: root}
}

func (fs *LocalFilesystem) path(key string) string {
	return filepath.Join(fs.root, key)
}

func (fs *LocalFilesystem) Read(key string) ([]byte, error) {
	path := fs.path(key)
	fp, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("asset %s %w", key, ErrNotFound)
		}
		return nil, fmt.Errorf("error opening asset %s: %w", key, err)
	}
	defer fp.Close()

	data, err := io.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("error reading asset %s: %w", key, err)
	}

	return data, nil
}

func (fs *LocalFilesystem) Write(key string, data []byte) error {
	path := fs.path(key)
	file, err := os.Create(path)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("asset %s exists: %w", key, ErrNotFound)
		}
		return fmt.Errorf("error opening asset %s: %w", key, err)
	}
	defer file.Close()
	_, err = file.Write(data)
	// todo: wrap errors
	return err
}
