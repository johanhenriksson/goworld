package engine

// Component is the general interface for scene object components.
type Component interface {
	Update(float32)
	Draw(DrawArgs)
}

// Draw a set of components.
func Draw(args DrawArgs, components ...Component) {
	for _, c := range components {
		c.Draw(args)
	}
}

// Update a set of components.
func Update(dt float32, components ...Component) {
	for _, c := range components {
		c.Update(dt)
	}
}
