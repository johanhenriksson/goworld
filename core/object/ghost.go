package object

func Ghost(object T) T {
	return &base{
		name:      "Ghost:" + object.Name(),
		enabled:   true,
		transform: object.Transform(),
	}
}
