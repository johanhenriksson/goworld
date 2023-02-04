package texture

import (
	"image"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Ref interface {
	Key() string
	Version() int

	// Load is called by texture caches and loaders, and should return the image data.
	// todo: This interface is a bit too simple as it does not allow us to pass
	//       formats, filters and aspects.
	Load() *image.RGBA
	Size() vec2.T
}

type pathRef struct {
	path string
	img  *image.RGBA
	size vec2.T
}

func PathRef(path string) Ref {
	return &pathRef{
		path: path,
		size: vec2.Zero,
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
	r.size = vec2.New(float32(r.img.Rect.Size().X), float32(r.img.Rect.Size().Y))
	return r.img
}

func (r *pathRef) Key() string  { return r.path }
func (r *pathRef) Version() int { return 1 }
func (r *pathRef) Size() vec2.T { return r.size }

type imageRef struct {
	name    string
	version int
	img     *image.RGBA
	size    vec2.T
}

func ImageRef(name string, version int, img *image.RGBA) Ref {
	return &imageRef{
		name:    name,
		version: version,
		img:     img,
		size:    vec2.New(float32(img.Rect.Size().X), float32(img.Rect.Size().Y)),
	}
}

func (r *imageRef) Key() string       { return r.name }
func (r *imageRef) Version() int      { return r.version }
func (r *imageRef) Load() *image.RGBA { return r.img }
func (r *imageRef) Size() vec2.T      { return r.size }

func (r *imageRef) String() string { return r.name }
