package object

type Array[T PropValue] struct {
	items []T
}

func NewArray[T PropValue]() Array[T] {
	return Array[T]{}
}

func (a *Array[T]) Length() int {
	return len(a.items)
}

func (a *Array[T]) Get(index int) T {
	return a.items[index]
}

func (a *Array[T]) GetAny(index int) any {
	return a.items[index]
}

func (a *Array[T]) Set(index int, value T) {
	a.items[index] = value
}

func (a *Array[T]) SetAny(index int, value any) {
	a.items[index] = value.(T)
}

func (a *Array[T]) Append(value T) {
	a.items = append(a.items, value)
}

func (a *Array[T]) AppendAny(value any) {
	a.items = append(a.items, value.(T))
}

func (a *Array[T]) Delete(index int) {
	a.items = append(a.items[:index], a.items[index+1:]...)
}

func (a *Array[T]) Serialize(enc Encoder) error {
	return enc.Encode(a.items)
}

func (a *Array[T]) Deserialize(dec Decoder) error {
	return dec.Decode(&a.items)
}
