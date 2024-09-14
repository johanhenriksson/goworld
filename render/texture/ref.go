package texture

import (
	"encoding/gob"

	"github.com/johanhenriksson/goworld/assets/fs"
	"github.com/johanhenriksson/goworld/render/image"
)

var Checker = PathRef("textures/uv_checker.png")

func init() {
	gob.Register(&pathRef{})
	gob.Register(Args{})
}

type Data struct {
	Args
	Image *image.Data
}

type pathRef struct {
	Path string
	Args Args

	data *Data
}

func PathRef(path string) *pathRef {
	return &pathRef{
		Path: path,
		Args: Args{
			Filter: FilterLinear,
			Wrap:   WrapRepeat,
		},
	}
}

func PathArgsRef(path string, args Args) *pathRef {
	return &pathRef{
		Path: path,
		Args: args,
	}
}

func (r *pathRef) Key() string  { return r.Path } // todo: this must include arguments somehow
func (r *pathRef) Version() int { return 1 }

func (r *pathRef) LoadTexture(assets fs.Filesystem) *Data {
	// caching
	// todo: move somewhere else where its easier to keep track of cached data
	if r.data != nil {
		return r.data
	}

	// load image
	img, err := image.LoadFile(assets, r.Path)
	if err != nil {
		panic(err)
	}

	r.data = &Data{
		Image: img,
		Args:  r.Args,
	}
	return r.data
}
