package ivec2

var Zero = T{}
var One = T{X: 1, Y: 1}
var UnitX = T{X: 1}
var UnitY = T{Y: 1}

type T struct {
	X int
	Y int
}

func New(x, y int) T {
	return T{
		X: x,
		Y: y,
	}
}

func (v T) Add(v2 T) T {
	return T{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v T) Sub(v2 T) T {
	return T{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}
