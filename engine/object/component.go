package object

import "github.com/johanhenriksson/goworld/math/vec3"

type Component interface {
	String() string
	Update(float32)

	Active() bool
	SetActive(bool) Component

	Parent() *T
	SetParent(*T) Component
	Collect(*Query)

	Forward() vec3.T
	Right() vec3.T
	Up() vec3.T

	Position() vec3.T
	SetPosition(vec3.T) Component

	Rotation() vec3.T
	SetRotation(vec3.T) Component

	Scale() vec3.T
	SetScale(vec3.T) Component
}
