package client

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/sprite"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/texture"
)

type Projectile struct {
	object.Object
	Sprite *sprite.Mesh

	velocity vec3.T
	lifetime float32
}

func NewProjectile(velocity vec3.T) *Projectile {
	return object.New("Projectile", &Projectile{
		Sprite: sprite.New(sprite.Args{
			Size: vec2.New(0.5, 0.5),
			Texture: texture.PathArgsRef("sprites/objects/barrel1.png", texture.Args{
				Filter: texture.FilterNearest,
			}),
		}),

		velocity: velocity,
		lifetime: 3,
	})
}

func (p *Projectile) Update(scene object.Component, dt float32) {
	p.Object.Update(scene, dt)

	p.lifetime -= dt
	if p.lifetime < 0 {
		p.Destroy()
		return
	}

	nextPos := p.Transform().Position().Add(p.velocity.Scaled(dt))
	p.Transform().SetPosition(nextPos)
}
