package object

import (
	"github.com/johanhenriksson/goworld/engine/transform"
)

type Component interface {
	String() string
	Update(float32)

	Active() bool
	SetActive(bool)

	Parent() T
	SetParent(T)
}

type T interface {
	String() string
	Update(float32)

	Active() bool
	SetActive(bool)

	Parent() T
	SetParent(T)

	Transform() transform.T

	Attach(...Component)

	Collect(*Query)

	Adopt(...T)
	Children() []T
}
