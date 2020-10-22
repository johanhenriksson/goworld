package engine

// Component is the general interface for scene components.
type Component interface {
	Name() string
	Update(float32)
	Collect(DrawPass, DrawArgs)
}

// Update a set of components.
func Update(dt float32, components ...Component) {
	for _, c := range components {
		if c == nil {
			continue
		}
		c.Update(dt)
	}
}

// Collect drawables from a set of components
func Collect(pass DrawPass, args DrawArgs, components ...Component) {
	for _, c := range components {
		if c == nil {
			continue
		}
		c.Collect(pass, args)
	}
}
