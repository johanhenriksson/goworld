package ivec2

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
