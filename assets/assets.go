package assets

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var ErrNotFound = fmt.Errorf("not found")

var Path string

const AssetPathConfig = "ASSET_PATH"

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	assetPath := "assets"
	if os.Getenv(AssetPathConfig) != "" {
		assetPath = os.Getenv(AssetPathConfig)
	}

	Path = FindFileInParents(assetPath, cwd)
}

func Open(key string) (io.ReadCloser, error) {
	fullpath := filepath.Join(Path, key)
	fp, err := os.Open(fullpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("asset %s %w", key, ErrNotFound)
		}
		return nil, fmt.Errorf("error opening asset %s: %w", key, err)
	}
	return fp, nil
}

func ReadAll(key string) ([]byte, error) {
	file, err := Open(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading asset %s: %w", key, err)
	}

	return data, nil
}

func Write(key string) (io.WriteCloser, error) {
	fullpath := filepath.Join(Path, key)
	file, err := os.Create(fullpath)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("asset %s exists: %w", key, ErrNotFound)
		}
		return nil, fmt.Errorf("error opening asset %s: %w", key, err)
	}

	return file, nil
}

func FindFileInParents(name, path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.Name() == name {
			return filepath.Join(path, name)
		}
	}
	return FindFileInParents(name, filepath.Dir(path))
}
