package object

// Returns all the children of an objects, both components and subgroups
func Children(object T) []T {
	if group, ok := object.(G); ok {
		return group.Children()
	}
	return nil
}

// Returns the subgroups attached to an object
func Subgroups(object T) []G {
	children := Children(object)
	groups := make([]G, 0, len(children))
	for _, child := range children {
		if group, ok := child.(G); ok {
			groups = append(groups, group)
		}
	}
	return groups
}

// Returns the components attached to an object
func Components(object T) []T {
	children := Children(object)
	components := make([]T, 0, len(children))
	for _, child := range children {
		_, group := child.(G)
		if !group {
			components = append(components, child)
		}
	}
	return components
}

// Attach an object to a parent object
// If the object already has a parent, it will be detached first.
func Attach(parent G, child T) {
	Detach(child)
	child.setParent(parent)
	parent.attach(child)

	if handler, ok := child.(ActivateHandler); ok {
		handler.OnActivate()
	}
}

// Detach an object from its parent object
// Does nothing if the given object has no parent.
func Detach(child T) {
	if child.Parent() == nil {
		return
	}
	child.Parent().detach(child)
	child.setParent(nil)
}

// Root returns the first ancestor of the given object
func Root(obj T) T {
	for obj.Parent() != nil {
		obj = obj.Parent()
	}
	return obj
}

// TODO: The Find functions should be rewritten as Queries.

// FindInParents finds the first object that implements the given interface
// in the root object or one of its ancestors.
func FindInParents[K T](root T) (K, bool) {
	if k, ok := root.(K); ok {
		return k, true
	}

	// consider components of direct ancestors
	for _, child := range Children(root) {
		if _, group := child.(G); group {
			continue
		}
		if k, ok := child.(K); ok {
			return k, true
		}
	}

	if root.Parent() != nil {
		return FindInParents[K](root.Parent())
	}
	var empty K
	return empty, false
}

// FindInChildren finds the first object that implements the given interface
// in the root object or one of its decendants.
func FindInChildren[K T](root T) (K, bool) {
	if k, ok := root.(K); ok {
		return k, true
	}
	if group, ok := root.(G); ok {
		// todo: rewrite as breadth-first
		for _, child := range group.Children() {
			if hit, ok := FindInChildren[K](child); ok {
				return hit, true
			}
		}
	}
	var empty K
	return empty, false
}

// FindInSiblings finds the first object that implements the given interface
// in the siblings of the given node.
func FindInSiblings[K T](self T) (K, bool) {
	var empty K
	if self.Parent() == nil {
		return empty, false
	}

	for _, child := range self.Parent().Children() {
		if child == self {
			continue
		}
		if hit, ok := child.(K); ok {
			return hit, true
		}
	}

	return empty, false
}

func FindAllInSiblings[K T](self T, callback func(K)) {
	if self.Parent() == nil {
		return
	}
	for _, child := range self.Parent().Children() {
		if child == self {
			continue
		}
		if hit, ok := child.(K); ok {
			callback(hit)
		}
	}
}
