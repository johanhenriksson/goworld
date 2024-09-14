package assets

import (
	"github.com/johanhenriksson/goworld/assets/fs"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/font"
	"github.com/johanhenriksson/goworld/render/texture"
)

type Texture interface {
	Asset

	// LoadImage is called by texture caches and loaders, and should return the image data.
	LoadTexture(assets fs.Filesystem) *texture.Data
}

var _ Texture = texture.PathRef("")
var _ Texture = (*font.Glyph)(nil)
var _ Texture = color.White
