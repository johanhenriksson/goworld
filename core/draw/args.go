package draw

import (
	"github.com/johanhenriksson/goworld/math/mat4"
)

// Args holds the arguments used to perform a draw pass.
// Includes the various transformation matrices and position of the camera.
type Args struct {
	Frame     int
	Time      float32
	Delta     float32
	Camera    Camera
	Transform mat4.T
}

// Apply the effects of a transform
func (d Args) Apply(t mat4.T) Args {
	d.Transform = d.Transform.Mul(&t)
	return d
}
