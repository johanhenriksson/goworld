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

type pathRef struct {
	Path string
	Args Args

	img *image.Data
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

func (r *pathRef) LoadImage(assets fs.Filesystem) *image.Data {
	if r.img != nil {
		return r.img
	}
	var err error
	r.img, err = image.LoadFile(assets, r.Path)
	if err != nil {
		panic(err)
	}
	return r.img
}

func (r *pathRef) TextureArgs() Args {
	return r.Args
}
