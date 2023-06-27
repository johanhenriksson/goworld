package object

type ghost struct {
	group
}

func Ghost(object T) G {
	return &ghost{
		group: group{
			base: base{
				id:      ID(),
				name:    "Ghost:" + object.Name(),
				enabled: true,
			},
			transform: object.Transform(),
		},
	}
}
