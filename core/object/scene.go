package object

import "log"

func init() {
	Register[*scene](deserializeScene)
}

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

func (s *scene) Serialize(enc Encoder) error {
	log.Println("serialize scene")
	return nil
}

func deserializeScene(dec Decoder) (Component, error) {
	log.Println("deserialize scene")
	return Scene(), nil
}
