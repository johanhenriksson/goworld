package render

type FloatArray []float32

func (a FloatArray) Elements() int {
	return len(a)
}

func (a FloatArray) Size() int {
	return 4
}

type UInt32Array []uint32

func (a UInt32Array) Elements() int {
	return len(a)
}

func (a UInt32Array) Size() int {
	return 4
}

type Int32Array []int32

func (a Int32Array) Elements() int {
	return len(a)
}

func (a Int32Array) Size() int {
	return 4
}
