package object

func Is[K any](c Component) bool {
	_, ok := c.(K)
	return ok
}
