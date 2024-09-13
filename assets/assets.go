package assets

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/johanhenriksson/goworld/assets/fs"
)

var ErrNotFound = fmt.Errorf("not found")

var FS fs.Filesystem

const AssetFolderEnv = "ASSET_PATH"

type Asset interface {
	Key() string
	Version() int
}

func init() {
	layeredFs := fs.NewLayered(BuiltinFilesystem)
	FS = layeredFs

	// look for a local asset path
	assetFolderName := "assets"
	if os.Getenv(AssetFolderEnv) != "" {
		assetFolderName = os.Getenv(AssetFolderEnv)
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if localAssetPath, err := FindFileInParents(assetFolderName, cwd); err == nil {
		log.Println("adding local file system layer rooted at", localAssetPath)
		layeredFs.Push(fs.NewLocal(localAssetPath))
	} else {
		log.Println("no local asset path found")
	}
}

func FindFileInParents(name, path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.Name() == name {
			return filepath.Join(path, name), nil
		}
	}
	parentPath := filepath.Dir(path)
	if parentPath == path {
		return "", ErrNotFound
	}
	return FindFileInParents(name, parentPath)
}
