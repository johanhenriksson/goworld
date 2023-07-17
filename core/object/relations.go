package object

// Returns all the children of an objects, both components and child objects
func Children(object Component) []Component {
	if group, ok := object.(Object); ok {
		return group.Children()
	}
	return nil
}

// Returns the child objects attached to an object
func Subgroups(object Component) []Object {
	children := Children(object)
	groups := make([]Object, 0, len(children))
	for _, child := range children {
		if group, ok := child.(Object); ok {
			groups = append(groups, group)
		}
	}
	return groups
}

// Returns the components attached to an object
func Components(object Component) []Component {
	children := Children(object)
	components := make([]Component, 0, len(children))
	for _, child := range children {
		_, group := child.(Object)
		if !group {
			components = append(components, child)
		}
	}
	return components
}

// Attach a child component/object to a parent object
// If the object already has a parent, it will be detached first.
func Attach(parent Object, child Component) {
	Detach(child)
	child.setParent(parent)
	parent.attach(child)
	activate(child)
}

// Detach a child component/object from its parent object
// Does nothing if the given object has no parent.
func Detach(child Component) {
	if child.Parent() == nil {
		return
	}
	deactivate(child)
	child.Parent().detach(child)
	child.setParent(nil)
}

func Enable(object Component) {
	object.setEnabled(true)
	activate(object)
}

func activate(object Component) {
	if !object.Enabled() {
		return
	}
	if object.Parent() == nil || !object.Parent().Active() {
		return
	}
	// activate if parent is active
	if wasActive := object.setActive(true); !wasActive {
		// enabled
		if handler, ok := object.(EnableHandler); ok {
			handler.OnEnable()
		}
	}
}

func Disable(object Component) {
	object.setEnabled(false)
	deactivate(object)
}

func deactivate(object Component) {
	if wasActive := object.setActive(false); wasActive {
		// disabled
		if handler, ok := object.(DisableHandler); ok {
			handler.OnDisable()
		}
	}
}

func Toggle(object Component, enabled bool) {
	if enabled {
		Enable(object)
	} else {
		Disable(object)
	}
}

// Root returns the first ancestor of the given component/object
func Root(obj Component) Component {
	for obj.Parent() != nil {
		obj = obj.Parent()
	}
	return obj
}

// Gets a reference to a component of type K on the same object as the component/object specified.
func Get[K Component](self Component) K {
	if hit, ok := self.(K); ok {
		return hit
	}
	var empty K
	group, ok := self.(Object)
	if !ok {
		group = self.Parent()
	}
	if group == nil {
		return empty
	}
	for _, child := range group.Children() {
		if child == self {
			continue
		}
		if hit, ok := child.(K); ok {
			return hit
		}
	}
	return empty
}

// Gets references to all components of type K on the same object as the component/object specified.
func GetAll[K Component](self Component) []K {
	group, ok := self.(Object)
	if !ok {
		group = self.Parent()
	}
	if group == nil {
		return nil
	}
	var results []K
	if hit, ok := group.(K); ok {
		results = append(results, hit)
	}
	for _, child := range group.Children() {
		if hit, ok := child.(K); ok {
			results = append(results, hit)
		}
	}
	return results
}

// Gets a reference to a component of type K on the same object as the component/object specified, or any parent of the object.
func GetInParents[K Component](self Component) K {
	var empty K
	group, ok := self.(Object)
	if !ok {
		group = self.Parent()
	}

	for group != nil {
		if hit, ok := group.(K); ok {
			return hit
		}
		for _, child := range group.Children() {
			if child == self {
				continue
			}
			if hit, ok := child.(K); ok {
				return hit
			}
		}
		group = group.Parent()
	}

	return empty
}

// Gets references to all components of type K on the same object as the component/object specified, or any parent of the object.
func GetAllInParents[K Component](self Component) []K {
	group, ok := self.(Object)
	if !ok {
		group = self.Parent()
	}
	var results []K
	for group != nil {
		if hit, ok := group.(K); ok {
			results = append(results, hit)
		}
		for _, child := range group.Children() {
			if child == self {
				continue
			}
			if hit, ok := child.(K); ok {
				results = append(results, hit)
			}
		}
		group = group.Parent()
	}
	return results
}

// Gets a reference to a component of type K on the same object as the component/object specified, or any child of the object.
func GetInChildren[K Component](self Component) K {
	var empty K
	group, ok := self.(Object)
	if !ok {
		group = self.Parent()
	}
	if group == nil {
		return empty
	}

	todo := []Object{group}

	for len(todo) > 0 {
		group = todo[0]
		todo = todo[1:]

		for _, child := range group.Children() {
			if child == self {
				continue
			}
			if hit, ok := child.(K); ok {
				return hit
			}
			if childgroup, ok := child.(Object); ok {
				todo = append(todo, childgroup)
			}
		}
	}

	return empty
}

// Gets references to all components of type K on the same object as the component/object specified, or any child of the object.
func GetAllInChildren[K Component](self Component) []K {
	group, ok := self.(Object)
	if !ok {
		group = self.Parent()
	}
	if group == nil {
		return nil
	}

	todo := []Object{group}
	var results []K

	for len(todo) > 0 {
		group = todo[0]
		todo = todo[1:]

		for _, child := range group.Children() {
			if child == self {
				continue
			}
			if hit, ok := child.(K); ok {
				results = append(results, hit)
			}
			if childgroup, ok := child.(Object); ok {
				todo = append(todo, childgroup)
			}
		}
	}

	return results
}
