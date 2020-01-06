package game

type OcclusionSpace interface {
	Free(x, y, z int) bool
}

type OcclusionData struct {
	data   []float32
	length int
	Size   int
}

func NewOcclusionData(size int) *OcclusionData {
	l := size * size * size
	return &OcclusionData{
		Size:   size,
		length: l,
		data:   make([]float32, l),
	}
}

func (o *OcclusionData) Get(x, y, z byte) byte {
	if x < 0 || y < 0 || z < 0 || int(x) >= o.Size || int(y) >= o.Size || int(z) >= o.Size {
		return 255
	}
	offset := int(z)*o.Size*o.Size + int(y)*o.Size + int(x)
	return byte(256 * o.data[offset])
}
func (o *OcclusionData) Set(x, y, z int, value float32) {
	if x < 0 || y < 0 || z < 0 || int(x) >= o.Size || int(y) >= o.Size || int(z) >= o.Size {
		return
	}
	offset := int(z)*o.Size*o.Size + int(y)*o.Size + int(x)

	o.data[offset] = value

}
