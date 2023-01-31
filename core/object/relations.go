package object

// Attach an object to a parent object
// If the object already has a parent, it will be detached first.
func Attach(parent, child T) {
	Detach(child)
	child.setParent(parent)
	parent.attach(child)
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
	// todo: rewrite as breadth-first
	for _, child := range root.Children() {
		if hit, ok := FindInChildren[K](child); ok {
			return hit, true
		}
	}
	var empty K
	return empty, false
}

// FindInSiblings finds the first object that implements the given interface
// in the siblings of the given root node.
func FindInSiblings[K T](root T) (K, bool) {
	if root.Parent() != nil {
		for _, child := range root.Parent().Children() {
			if child == root {
				continue
			}
			if hit, ok := child.(K); ok {
				return hit, true
			}
		}
	}
	var empty K
	return empty, false
}
