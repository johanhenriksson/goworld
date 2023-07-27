package physics

import (
	"github.com/johanhenriksson/goworld/math/vec3"
)

type RaycastHit struct {
	Shape  Shape
	Point  vec3.T
	Normal vec3.T
}

func (w *World) Raycast(from, to vec3.T, mask Mask) (hit RaycastHit, exists bool) {
	result, didHit := world_raycast(w.handle, from, to, mask)
	if didHit {
		exists = true
		hit = RaycastHit{
			Shape:  restoreShape(result.shape),
			Point:  result.point,
			Normal: result.normal,
		}
	}
	return
}
