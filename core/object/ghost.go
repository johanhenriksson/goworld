package object

func Ghost(object T) T {
	return &base{
		id:        ID(),
		name:      "Ghost:" + object.Name(),
		enabled:   true,
		transform: object.Transform(),
	}
}
