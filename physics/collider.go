package physics

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Collider interface {
	OnCollision(other Collider, contact Contact)
	AttachToBody(rb *RigidBody)
}

/* Collision contact point */
type Contact struct {
	Position mgl.Vec3
	Normal   mgl.Vec3
	Depth    float32
}

type CollisionCallback func(other Collider, contact Contact)
