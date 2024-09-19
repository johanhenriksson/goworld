package assets

import (
	"github.com/johanhenriksson/goworld/assets/fs"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/font"
	"github.com/johanhenriksson/goworld/render/texture"
)

type Texture interface {
	Asset

	// LoadTexture is called by texture caches and loaders, and should return the texture data.
	// Its unfortunate that this method cant return the texture itself directly, since it
	// requires access to a graphics queue. The texture upload logic must be centralized somewhere.
	// Techincally not different for meshes? hmm
	// Should there be a concept on a cpu-side texture? Similar to vertex.Mesh?
	LoadTexture(assets fs.Filesystem) *texture.Data
}

var _ Texture = texture.PathRef("")
var _ Texture = (*font.Glyph)(nil)
var _ Texture = color.White
