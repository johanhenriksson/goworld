package object

type Dict[K comparable, V any] struct {
	items map[K]V
}

func NewDict[K comparable, V any]() Dict[K, V] {
	return Dict[K, V]{
		items: map[K]V{},
	}
}

func (d *Dict[K, V]) Get(key K) (V, bool) {
	v, ok := d.items[key]
	return v, ok
}

func (d *Dict[K, V]) GetAny(key K) (any, bool) {
	v, ok := d.items[key]
	return v, ok
}

func (d *Dict[K, V]) Set(key K, value V) {
	d.items[key] = value
}

func (d *Dict[K, V]) SetAny(key K, value any) {
	d.items[key] = value.(V)
}

func (d *Dict[K, V]) Delete(key K) {
	delete(d.items, key)
}

func (d *Dict[K, V]) Serialize(enc Encoder) error {
	return enc.Encode(d.items)
}

func (d *Dict[K, V]) Deserialize(dec Decoder) error {
	return dec.Decode(&d.items)
}
