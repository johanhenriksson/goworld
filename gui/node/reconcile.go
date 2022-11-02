package node

import "fmt"

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
		if _, exists := previous[child.Key()]; exists {
			panic(fmt.Errorf("duplicate key %s in children of %s", child.Key(), target.Key()))
		}
		previous[child.Key()] = child
	}

	children := new.Children()
	for idx, child := range children {
		// todo: handle nil children

		if existing, ok := previous[child.Key()]; ok {
			// since each key can only appear once, we can remove the child from the mapping
			// this prevents the old element from being destroyed later. if its no longer needed,
			// it will be destroyed in the reconciliation below.
			delete(previous, child.Key())

			// recursively reconcile child node
			children[idx] = Reconcile(existing, child)
		} else {
			// this key did not exist previously, so it must be a new element
		}
	}

	// replace reconciled children
	target.SetChildren(children)

	// at this point, any child remaining in the previous map is no longer part of the tree
	// and should be destroyed
	for _, child := range previous {
		child.Destroy()
	}

	return target
}
