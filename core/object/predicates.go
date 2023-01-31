package object

func Is[K any](c T) bool {
	_, ok := c.(K)
	return ok
}
