package assets

import (
	"fmt"
	"os"
	"path/filepath"
)

var ErrNotFound = fmt.Errorf("not found")

var Path string
var server Server

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

	server = NewLayeredFilesystem(
		NewLocalFilesystem(Path),
		BuiltinFilesystem,
	)
}

func Read(key string) ([]byte, error) {
	return server.Read(key)
}

func Write(key string, data []byte) error {
	return server.Write(key, data)
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
