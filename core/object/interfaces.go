package object

import (
	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/core/transform"
)

type Updatable interface {
	// Update the object. Called on every frame.
	Update(float32)

	// Active indicates whether the object is currently enabled or not.
	Active() bool

	// SetActive enables or disables the object
	SetActive(bool)
}

type Transformed interface {
	Transform() transform.T
}

type Component interface {
	Updatable
	Transformed

	// String returns the name of the component
	Name() string

	// Object returns the parent object the component is attached to, or nil
	Object() T

	// SetObject attaches the component to an object
	SetObject(T)
}

// T is a object in the scene hierarchy
type T interface {
	Updatable
	Transformed
	input.Handler

	// String is the name of the object
	Name() string

	Parent() T
	SetParent(T)

	// Adopt a child object
	Adopt(...T)
	Children() []T

	Attach(...Component)
	Collect(*Query)
}
