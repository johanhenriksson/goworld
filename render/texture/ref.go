package texture

import (
	"github.com/johanhenriksson/goworld/render/image"
)

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
}

func PathRef(path string) Ref {
	return &pathRef{
		path: path,
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

	return Args{}
}
