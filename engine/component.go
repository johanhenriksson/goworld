package engine

// Component is the general interface for scene object components.
type Component interface {
	// Base() *Object
	Update(float32)
	Draw(DrawArgs)
}
