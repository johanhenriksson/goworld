package texture

import (
	"image"

	"github.com/johanhenriksson/goworld/assets"
)

type Ref interface {
	Key() string
	Version() int

	// Load is called by texture caches and loaders, and should return the image data.
	// todo: This interface is a bit too simple as it does not allow us to pass
	//       formats, filters and aspects.
	Load() *image.RGBA
}

type pathRef struct {
	path string
	img  *image.RGBA
}

func PathRef(path string) Ref {
	return &pathRef{
		path: path,
	}
}

func (r *pathRef) Load() *image.RGBA {
	if r.img != nil {
		return r.img
	}
	var err error
	r.img, err = assets.GetImage(r.path)
	if err != nil {
		panic(err)
	}
	return r.img
}

func (r *pathRef) Key() string  { return r.path }
func (r *pathRef) Version() int { return 1 }

type imageRef struct {
	name    string
	version int
	img     *image.RGBA
}

func ImageRef(name string, version int, img *image.RGBA) Ref {
	return &imageRef{
		name:    name,
		version: version,
		img:     img,
	}
}

func (r *imageRef) Key() string       { return r.name }
func (r *imageRef) Version() int      { return r.version }
func (r *imageRef) Load() *image.RGBA { return r.img }

func (r *imageRef) String() string { return r.name }
