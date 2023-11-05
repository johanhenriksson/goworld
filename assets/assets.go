package assets

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var vfs fs.FS
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
	vfs = os.DirFS(Path)
}

func Open(fileName string) (fs.File, error) {
	return vfs.Open(fileName)
}

func ReadAll(fileName string) ([]byte, error) {
	file, err := Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", fileName, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", fileName, err)
	}

	return data, nil
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
