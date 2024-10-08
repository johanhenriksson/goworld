package uniform

import (
	"structs"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render/device"
)

const MaxTextures = 16

type TextureId uint32
type TextureIds [MaxTextures]TextureId

type Object struct {
	_ structs.HostLayout

	Model    mat4.T
	Textures TextureIds

	Vertices device.Address
	Indices  device.Address
}
