package dimension

type T interface {
	Resolve(float32) float32
	Fixed() bool
	Auto() bool
}

type D struct {
	Value T
	Grow  bool
}

// flex 1 1 20px

type Fixed float32

func (f Fixed) Resolve(float32) float32 { return float32(f) }
func (f Fixed) Fixed() bool             { return true }
func (f Fixed) Auto() bool              { return false }

func Auto() T {
	return auto{}
}

type auto struct{}

func (a auto) Resolve(parent float32) float32 { return parent }
func (a auto) Fixed() bool                    { return false }
func (a auto) Auto() bool                     { return true }

type Percent float32

func (p Percent) Resolve(parent float32) float32 { return 0.01 * float32(p) }
func (f Percent) Fixed() bool                    { return false }
func (f Percent) Auto() bool                     { return false }
