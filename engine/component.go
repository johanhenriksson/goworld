package engine

// Component is the general interface for scene object components.
type Component interface {
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

func Collect(pass DrawPass, args DrawArgs, components ...Component) {
	for _, c := range components {
		if c == nil {
			continue
		}
		c.Collect(pass, args)
	}
}
