package object

type scene struct {
	Object
}

func Scene() Object {
	return &scene{
		Object: Empty("Scene"),
	}
}

func (s *scene) Active() bool {
	return true
}

func (s *scene) setActive(bool) bool {
	return true
}
