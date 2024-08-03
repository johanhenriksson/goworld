package spline

import (
	"encoding/gob"

	"github.com/johanhenriksson/goworld/math/vec2"
)

type Linear struct {
	Points []vec2.T
}

func init() {
	gob.Register(Linear{})
}

func NewLinear(points ...vec2.T) Linear {
	return Linear{
		Points: points,
	}
}

func (l Linear) Eval(t float32) float32 {
	if len(l.Points) == 0 {
		return 0
	}
	if len(l.Points) == 1 {
		return l.Points[0].Y
	}
	if t <= l.Points[0].X {
		return l.Points[0].Y
	}
	for i := 1; i < len(l.Points); i++ {
		if t < l.Points[i].X {
			a := l.Points[i-1]
			b := l.Points[i]
			return a.Y + (t-a.X)*(b.Y-a.Y)/(b.X-a.X)
		}
	}
	return l.Points[len(l.Points)-1].Y
}
