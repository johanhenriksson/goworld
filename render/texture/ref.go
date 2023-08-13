package texture

import (
	"github.com/johanhenriksson/goworld/render/image"
)

var Checker = PathRef("textures/uv_checker.png")

type Ref interface {
	Key() string
	Version() int

	// ImageData is called by texture caches and loaders, and should return the image data.
	// todo: This interface is a bit too simple as it does not allow us to pass
	//       formats, filters and aspects.
	ImageData() *image.Data
	TextureArgs() Args
}

type pathRef struct {
	path string
	img  *image.Data
	args Args
}

func PathRef(path string) Ref {
	return &pathRef{
		path: path,
		args: Args{
			Filter: FilterLinear,
			Wrap:   WrapRepeat,
		},
	}
}

func PathArgsRef(path string, args Args) Ref {
	return &pathRef{
		path: path,
		args: args,
	}
}

func (r *pathRef) Key() string  { return r.path }
func (r *pathRef) Version() int { return 1 }

func (r *pathRef) ImageData() *image.Data {
	if r.img != nil {
		return r.img
	}
	var err error
	r.img, err = image.LoadFile(r.path)
	if err != nil {
		panic(err)
	}
	return r.img
}

func (r *pathRef) TextureArgs() Args {
	return r.args
}
