package engine

// Component is the general interface for scene object components.
type Component interface {
	// Base() *Object
	Update(float32)
	Draw(DrawArgs)
}

func Draw(args DrawArgs, components ...Component) {
	for _, c := range components {
		c.Draw(args)
	}
}

func Update(dt float32, components ...Component) {
	for _, c := range components {
		c.Update(dt)
	}
}
