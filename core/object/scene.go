package object

func init() {
	Register[*scene](TypeInfo{
		Name: "Scene",
		Create: func(pool Pool) (Component, error) {
			return Scene(pool), nil
		},
	})
}

type SceneFunc func(Pool, Object)

type scene struct {
	Object
}

func Scene(pool Pool, funcs ...SceneFunc) Object {
	s := &scene{
		Object: Empty(pool, "Scene"),
	}
	for _, f := range funcs {
		f(pool, s)
	}
	return s
}

func (s *scene) Active() bool {
	return true
}

func (s *scene) setActive(bool) bool {
	return true
}
