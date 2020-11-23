package object

import "github.com/johanhenriksson/goworld/math/vec3"

type Component interface {
	String() string
	Update(float32)

	Active() bool
	SetActive(bool)

	Parent() *T
	SetParent(*T)
	Collect(*Query)

	Forward() vec3.T
	Right() vec3.T
	Up() vec3.T

	Position() vec3.T
	SetPosition(vec3.T)

	Rotation() vec3.T
	SetRotation(vec3.T)

	Scale() vec3.T
	SetScale(vec3.T)
}
