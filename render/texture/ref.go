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

func PathRef(path string) Ref {
	img, err := assets.GetImage(path)
	if err != nil {
		panic(err)
	}
	return ImageRef(path, 1, img)
}

type image_ref struct {
	name    string
	version int
	img     *image.RGBA
	size    vec2.T
}

func ImageRef(name string, version int, img *image.RGBA) Ref {
	return &image_ref{
		name:    name,
		version: version,
		img:     img,
		size:    vec2.New(float32(img.Rect.Size().X), float32(img.Rect.Size().Y)),
	}
}

func (r *image_ref) Key() string       { return r.name }
func (r *image_ref) Version() int      { return r.version }
func (r *image_ref) Load() *image.RGBA { return r.img }
func (r *image_ref) Size() vec2.T      { return r.size }
