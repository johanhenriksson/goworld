package style

import ()

type Font struct {
	Name string
	Size int
}

func (f Font) ApplyFont(fw FontWidget) {
	fw.SetFont(f)
}
