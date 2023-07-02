package object

import "github.com/johanhenriksson/goworld/core/transform"

type ghost struct {
	group
	target T
}

func Ghost(object T) G {
	return &ghost{
		group: group{
			base: base{
				id:      ID(),
				name:    "Ghost:" + object.Name(),
				enabled: true,
			},
		},
		target: object,
	}
}

func (g *ghost) Transform() transform.T {
	return g.target.Transform()
}
