package physics

type Mask uint32

const All = Mask(0xFFFFFFFF)
const None = Mask(0)

func Layers(layers ...Mask) Mask {
	mask := None
	for _, layer := range layers {
		mask = mask | layer
	}
	return mask
}
