package physics

/*
#cgo CXXFLAGS: -std=c++11 -I/usr/local/include/bullet
#cgo CFLAGS: -I/usr/local/include/bullet
#cgo LDFLAGS: -lstdc++ -L/usr/local/lib -lBulletDynamics -lBulletCollision -lLinearMath -lBullet3Common
#include "bullet.h"
*/
import "C"
import (
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

//export GoDebugCallback
func GoDebugCallback(start_x, start_y, start_z, end_x, end_y, end_z, color_r, color_g, color_b C.float) {
	start := vec3.New(float32(start_x), float32(start_y), float32(start_z))
	end := vec3.New(float32(end_x), float32(end_y), float32(end_z))
	color := color.RGB(float32(color_r), float32(color_g), float32(color_b))
	// fmt.Printf("DrawDebugLine: %s -> %s (%s)\n", start, end, color)

	lines.Debug.Add(start, end, color)
}
