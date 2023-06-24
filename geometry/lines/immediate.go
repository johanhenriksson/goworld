package lines

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

var Debug = New(Args{
	Lines: []Line{
		{
			Start: vec3.New(0, 0, 0),
			End:   vec3.New(0, 100, 0),
			Color: color.Yellow,
		},
	},
})
