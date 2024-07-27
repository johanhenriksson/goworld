package uniform

import "github.com/johanhenriksson/goworld/math/mat4"

const MaxTextures = 16

type TextureId uint32
type TextureIds [MaxTextures]TextureId

type Object struct {
	Model    mat4.T
	Textures TextureIds
}
