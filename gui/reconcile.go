package gui

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/johanhenriksson/goworld/gui/rect"
	"github.com/johanhenriksson/goworld/gui/widget"
)

func Reconcile(src, dst widget.T) bool {
	return reconcile(src, dst, 0)
}

func reconcile(src, dst widget.T, depth int) bool {
	indent := strings.Repeat("  ", depth)

	// check type
	srcType, dstType := reflect.TypeOf(src), reflect.TypeOf(dst)
	fmt.Printf("%s RECONCILE %s:%s -> %s:%s\n", indent, src.Key(), reflect.TypeOf(src), dst.Key(), reflect.TypeOf(dst))
	if srcType != dstType {
		fmt.Println(indent, "! types differ:", srcType, "!=", dstType)

		// clean up old element
		src.Destroy()

		fmt.Println(indent, "FAIL")
		return false
	}

	if src.Key() != dst.Key() {
		// if the keys dont match, reconcilation is not considered
		// at this point we can discard all the elements in the old (sub)tree
		fmt.Println(indent, "! keys differ:", src.Key(), "!=", dst.Key())

		// clean up old element
		src.Destroy()

		fmt.Println(indent, "FAIL")
		return false
	}

	// compare props
	srcprops := src.Props()
	dstprops := dst.Props()
	if !reflect.DeepEqual(srcprops, dstprops) {
		// props are NOT equal
		// we need to update them
		// this will possibly cause a reflow event
		fmt.Printf("%s ~ props differ: %+v vs %+v\n", indent, srcprops, dstprops)
		src.Update(dst.Props())
	}

	// reconcile children
	if dstRect, ok := dst.(rect.T); ok {
		srcRect := src.(rect.T)

		children := dstRect.Children()

		// create a key mapping for the existing child nodes
		previous := map[string]widget.T{}
		for _, child := range srcRect.Children() {
			previous[child.Key()] = child
		}

		for idx, child := range children {
			if existing, ok := previous[child.Key()]; ok {
				// since each key can only appear once, we can remove the child from the mapping
				delete(previous, child.Key())

				if reconcile(child, existing, depth+1) {
					// subtree reconciliation was successful!
					// replace the new child with the existing one
					children[idx] = existing
					fmt.Println(indent, "* reuse", existing.Key(), "at index", idx)
				} else {
					// unable to reconcile child!
					// destroy the old one.
					fmt.Println(indent, "! recreate", child.Key(), "at index", idx)
					existing.Destroy()
				}
			} else {
				// this key did not exist previously, so it must be a new element
				fmt.Println(indent, "! create", child.Key(), "at index", idx)
			}
		}

		// replace source children
		srcRect.SetChildren(children)

		// at this point, any child left in the previous map should be destroyed
		for _, child := range previous {
			fmt.Println(indent, "! removed", child.Key())
			child.Destroy()
		}

		// clear the child list of the dst rect
		// we have manually destroyed the ones that we wont reuse
		// if we dont, our reused children will be destroyed
		dstRect.SetChildren([]widget.T{})
	}

	// destroy dst
	dst.Destroy()

	fmt.Println(indent, "OK")
	return true
}
