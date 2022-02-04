package node

func Reconcile(target, new T) T {
	// no source tree - just go with the new one
	if target == nil {
		target = new
	}

	// element types are different - so we obviously can not reconcile
	// if the keys dont match, reconcilation is not considered
	if target.Type() != new.Type() || target.Key() != new.Key() {
		target.Destroy()
		target = new
	}

	// update props
	target.Update(new.Props())

	// expand new node to look at its children
	// use the existing hook state
	new.Render(target.Hooks())

	// create a key mapping for the existing child nodes
	// this allows us to reuse nodes and keep track of deletions
	previous := map[string]T{}
	for _, child := range target.Children() {
		// todo: check for duplicate keys
		previous[child.Key()] = child
	}

	children := new.Children()
	for idx, child := range children {
		// todo: handle nil children

		if existing, ok := previous[child.Key()]; ok {
			// since each key can only appear once, we can remove the child from the mapping
			delete(previous, child.Key())

			// recursively reconcile child node
			children[idx] = Reconcile(existing, child)
		} else {
			// this key did not exist previously, so it must be a new element
		}
	}

	// replace source children
	target.SetChildren(children)

	// at this point, any child left in the previous map should be destroyed
	for _, child := range previous {
		child.Destroy()
	}

	return target
}
