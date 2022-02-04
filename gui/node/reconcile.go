package node

func Reconcile(src, dst T) T {
	// no source tree - just go with the new one
	if src == nil {
		return dst
	}

	// element types are different - so we obviously can not reconcile
	// if the keys dont match, reconcilation is not considered
	if src.Type() != dst.Type() || src.Key() != dst.Key() {
		src.Destroy()
		return dst
	}

	// we can reuse the existing element!
	// update source props
	src.Update(dst.Props())

	// reconcile children - render the node so we can inspect them
	dst.Render()
	children := dst.Children()

	// create a key mapping for the existing child nodes
	// this allows us to reuse nodes and keep track of deletions
	previous := map[string]T{}
	for _, child := range src.Children() {
		// todo: check for duplicate keys
		previous[child.Key()] = child
	}

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
	src.SetChildren(children)

	// at this point, any child left in the previous map should be destroyed
	for _, child := range previous {
		child.Destroy()
	}

	return src
}
