package object

import "github.com/johanhenriksson/goworld/core/transform"

type ghost struct {
	group
	target Component
}

func Ghost(object Component) G {
	return &ghost{
		group: group{
			base: base{
				id:      ID(),
				name:    "Ghost:" + object.Name(),
				enabled: true,
			},
			transform: transform.Identity(),
		},
		target: object,
	}
}

func (g *ghost) Transform() transform.T {
	return g.target.Transform()
}
