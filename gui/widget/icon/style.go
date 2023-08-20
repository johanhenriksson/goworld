package icon

import (
	. "github.com/johanhenriksson/goworld/gui/style"
)

type Style struct {
	Hover Hover
	Size  int
	Color ColorProp
}

type Hover struct {
	Color ColorProp
}
