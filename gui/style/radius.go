package style

type RadiusProp interface {
	ApplyRadius(RadiusWidget)
}

type RadiusWidget interface {
	SetRadius(px float32)
}

func (p Px) ApplyRadius(w RadiusWidget) {
	w.SetRadius(float32(p))
}
