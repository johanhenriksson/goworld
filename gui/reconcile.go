package gui

import (
	"reflect"

	"github.com/johanhenriksson/goworld/gui/widget"
)

func Reconcile(src, dst widget.T) bool {
	return reconcile(src, dst, 0)
}

func reconcile(src, dst widget.T, depth int) bool {
	if src == nil {
		return false
	}

	// compare element type
	srcType, dstType := reflect.TypeOf(src), reflect.TypeOf(dst)
	if srcType != dstType {
		// element types are different - so we obviously can not reconcile

		// clean up old element
		src.Destroy()

		return false
	}

	if src.Key() != dst.Key() {
		// if the keys dont match, reconcilation is not considered
		// at this point we can discard all the elements in the old (sub)tree

		// clean up old element
		src.Destroy()

		return false
	}

	// compare props
	srcprops := src.Props()
	dstprops := dst.Props()
	if !reflect.DeepEqual(srcprops, dstprops) {
		// props are NOT equal
		// we need to update them
		// this will possibly cause a reflow event
		src.Update(dst.Props())
	}

	// reconcile children - if src and dst are Rects
	children := dst.Children()

	// create a key mapping for the existing child nodes
	previous := map[string]widget.T{}
	for _, child := range src.Children() {
		previous[child.Key()] = child
	}

	for idx, child := range children {
		// todo: handle nil children

		if existing, ok := previous[child.Key()]; ok {
			// since each key can only appear once, we can remove the child from the mapping
			delete(previous, child.Key())

			if reconcile(existing, child, depth+1) {
				// subtree reconciliation was successful!
				// replace the new child with the existing one
				children[idx] = existing
			} else {
				// unable to reconcile child!
				// destroy the old one.
				existing.Destroy()
			}
		}
		// this key did not exist previously, so it must be a new element
	}

	// replace source children
	src.SetChildren(children)

	// at this point, any child left in the previous map should be destroyed
	for _, child := range previous {
		child.Destroy()
	}

	// clear the child list of the dst rect
	// we have manually destroyed the ones that we wont reuse
	// if we dont, our reused children will be destroyed
	dst.SetChildren(nil)

	// destroy dst
	dst.Destroy()

	return true
}
