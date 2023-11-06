package client

import (
	"fmt"
	"time"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/sprite"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/texture"
)

const FacingToward = 0
const FacingAway = 1
const FacingLeft = 2
const FacingRight = 3

type CharacterSprite struct {
	object.Object
	Sprite *sprite.Mesh
	Index  int
	Frame  int
	Facing int

	animating bool
	frames    int
	frameTime float32
	nextFrame float32
}

func NewCharacterSprite(id int) *CharacterSprite {
	ft := float32(time.Duration(time.Second / 10).Seconds())
	return object.New("CharacterSprite", &CharacterSprite{
		Sprite: sprite.New(sprite.Args{
			Size: vec2.New(2, 2),
		}),
		Index: id,

		frames:    3,
		frameTime: ft,
		nextFrame: ft,
	})
}

func (c *CharacterSprite) Update(scene object.Component, dt float32) {
	cam := object.GetInChildren[*camera.Camera](scene)
	if cam == nil {
		return
	}

	c.nextFrame -= dt
	if c.nextFrame < 0 {
		c.nextFrame += c.frameTime
		if c.animating {
			c.Frame = (c.Frame + 1) % 2
		}
	}

	spriteFwd := c.Transform().Forward()
	spriteFwd.Y = 0
	spriteFwd.Normalize()
	camFwd := cam.Transform().ProjectDir(vec3.Forward)
	camFwd.Y = 0
	camFwd.Normalize()
	forward := vec3.Dot(spriteFwd, camFwd)

	const limit = float32(0.707) // 45deg
	if forward > limit {
		c.Facing = FacingToward
	} else if forward < -limit {
		c.Facing = FacingAway
	} else {
		// left or right
		camRight := cam.Transform().ProjectDir(vec3.Right)
		camRight.Y = 0
		camRight.Normalize()
		right := vec3.Dot(spriteFwd, camRight)

		if right > 0 {
			c.Facing = FacingRight
		} else {
			c.Facing = FacingLeft
		}
	}
	frame := c.Facing*c.frames + c.Frame

	c.Sprite.Sprite.Set(texture.PathArgsRef(fmt.Sprintf("sprites/sprite_%d_%d.png", c.Index, frame), texture.Args{
		Filter: texture.FilterNearest,
	}))
}
