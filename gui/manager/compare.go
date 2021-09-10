package manager

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/gui/widget"
)

func compare(src, dst widget.T) {
	// compare props
	srcprops := src.Props()
	dstprops := dst.Props()
	fmt.Println(srcprops)
	fmt.Println(dstprops)

	// aptr := reflect.ValueOf(srcprops).Elem().Interface()
	// bptr := reflect.ValueOf(dstprops).Elem().Interface()

	if reflect.DeepEqual(srcprops, dstprops) {
		fmt.Println("props are equal")
	} else {
		fmt.Println("props are NOT equal")
	}

	// compare children
}
