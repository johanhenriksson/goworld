package draw

import "github.com/johanhenriksson/goworld/math/vec2"

type Viewport struct {
	Width  int
	Height int
	Scale  float32
}

func (s Viewport) Aspect() float32 {
	return float32(s.Width) / float32(s.Height)
}

func (s Viewport) Size() vec2.T {
	return vec2.NewI(s.Width, s.Height)
}

func (s Viewport) NormalizeCursor(cursor vec2.T) vec2.T {
	return cursor.Div(s.Size()).Sub(vec2.New(0.5, 0.5)).Scaled(2)
}
