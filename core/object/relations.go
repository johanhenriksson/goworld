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

// Gets a reference to a component of type K in the same group as the object specified.
func Get[K T](self T) K {
	var empty K
	group, ok := self.(G)
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

// Gets references to all components of type K in the same group as the object specified.
func GetAll[K T](self T) []K {
	group, ok := self.(G)
	if !ok {
		group = self.Parent()
	}
	if group == nil {
		return nil
	}
	var results []K
	for _, child := range group.Children() {
		if hit, ok := child.(K); ok {
			results = append(results, hit)
		}
	}
	return results
}

// Gets a reference to a component of type K in the same group as the component specified, or any parent of the group.
func GetInParents[K T](self T) K {
	var empty K
	group, ok := self.(G)
	if !ok {
		group = self.Parent()
	}

	for group != nil {
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

// Gets references to all components of type K in the same group as the object specified, or any parent of the group.
func GetAllInParents[K T](self T) []K {
	group, ok := self.(G)
	if !ok {
		group = self.Parent()
	}
	var results []K
	for group != nil {
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

// Gets a reference to a component of type K in the same group as the object specified, or any child of the group.
func GetInChildren[K T](self T) K {
	var empty K
	group, ok := self.(G)
	if !ok {
		group = self.Parent()
	}
	if group == nil {
		return empty
	}

	todo := []G{group}

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
			if childgroup, ok := child.(G); ok {
				todo = append(todo, childgroup)
			}
		}
	}

	return empty
}

// Gets references to all components of type K in the same group as the object specified, or any child of the group.
func GetAllInChildren[K T](self T) []K {
	group, ok := self.(G)
	if !ok {
		group = self.Parent()
	}
	if group == nil {
		return nil
	}

	todo := []G{group}
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
			if childgroup, ok := child.(G); ok {
				todo = append(todo, childgroup)
			}
		}
	}

	return results
}
