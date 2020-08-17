package ui

func RowLayout(c Component, style Style, aw, ah float32) (float32, float32) {
	pad := style.Float("padding", 0)
	spacing := style.Float("spacing", 0)

	dw := float32(0)
	dh := float32(0)
	for _, child := range c.Children() {
		cx, cy := dw+pad, pad
		cdw, cdh := child.DesiredSize(aw-dw-2*pad, ah-2*pad)
		child.SetPosition(cx, cy)
		dw += cdw + spacing
		if cdh > dh {
			dh = cdh
		}
	}
	dw += 2*pad - spacing
	dh += 2 * pad

	c.SetSize(dw, dh)
	return dw, dh
}

func ColumnLayout(c Component, style Style, aw, ah float32) (float32, float32) {
	pad := style.Float("padding", 0)
	spacing := style.Float("spacing", 0)

	dw := float32(0)
	dh := float32(0)
	for _, child := range c.Children() {
		cx, cy := float32(pad), dh+pad
		cdw, cdh := child.DesiredSize(aw-2*pad, ah-dh-2*pad)
		child.SetPosition(cx, cy)
		dh += cdh + spacing
		if cdw > dw {
			dw = cdw
		}
	}
	dh += 2*pad - spacing
	dw += 2 * pad

	c.SetSize(dw, dh)
	return dw, dh
}
