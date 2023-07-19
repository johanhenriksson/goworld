package object

import "github.com/johanhenriksson/goworld/core/transform"

type ghost struct {
	object
	target Component
}

func Ghost(obj Component) Object {
	return &ghost{
		object: object{
			component: component{
				id:      ID(),
				name:    "Ghost:" + obj.Name(),
				enabled: true,
			},
			transform: transform.Identity(),
		},
		target: obj,
	}
}

func (g *ghost) Transform() transform.T {
	return g.target.Transform()
}
