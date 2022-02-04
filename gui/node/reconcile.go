package node

func Reconcile(src, dst T) T {
	if src == nil {
		return dst
	}

	// compare element type
	if src.Type() != dst.Type() {
		// element types are different - so we obviously can not reconcile

		// clean up old element
		src.Destroy()

		return dst
	}

	if src.Key() != dst.Key() {
		// if the keys dont match, reconcilation is not considered
		// at this point we can discard all the elements in the old (sub)tree

		// clean up old element
		src.Destroy()

		return dst
	}

	// compare props
	src.Update(dst.Props())

	// reconcile children
	dst.Render()
	children := dst.Children()

	// create a key mapping for the existing child nodes
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

			children[idx] = Reconcile(existing, child)
		} else {
			// this key did not exist previously, so it must be a new element
		}
	}

	// replace source children
	src.SetChildren(children)

	// at this point, any child left in the previous map should be destroyed
	for _, child := range previous {
		// unmount
		child.Destroy()
	}

	return src
}
