package assets

import (
	"fmt"
	"os"
	"path/filepath"
)

var ErrNotFound = fmt.Errorf("not found")

var Path string
var assetFs Filesystem

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

	assetFs = NewLayeredFilesystem(
		NewLocalFilesystem(Path),
		BuiltinFilesystem,
	)
}

func Read(key string) ([]byte, error) {
	return assetFs.Read(key)
}

func Write(key string, data []byte) error {
	return assetFs.Write(key, data)
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
	parentPath := filepath.Dir(path)
	if parentPath == path {
		return ""
	}
	return FindFileInParents(name, parentPath)
}
