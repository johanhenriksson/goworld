package dimension

type T interface {
	Resolve(float32) float32
}

type Fixed float32

func (f Fixed) Resolve(float32) float32 { return float32(f) }

func Auto() T {
	return auto{}
}

type auto struct{}

func (a auto) Resolve(parent float32) float32 { return parent }

type Percent float32

func (p Percent) Resolve(parent float32) float32 { return 0.01 * float32(p) * parent }
