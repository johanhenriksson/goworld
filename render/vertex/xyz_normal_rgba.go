package vertex

import (
	"github.com/johanhenriksson/goworld/math/byte4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// FPNC - Position, Normal, Color
type FPNC struct {
	P vec3.T  `vtx:"position,float,3"`
	N vec3.T  `vtx:"normal,float,3"`
	C byte4.T `vtx:"color,byte,4,normalize"`
}

// FPNT - Position, Normal, Texcoords
type FPNT struct {
}
