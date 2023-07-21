package physics

/*
#cgo CXXFLAGS: -std=c++11 -I/usr/local/include/bullet
#cgo CFLAGS: -I/usr/local/include/bullet
#cgo LDFLAGS: -lstdc++ -L/usr/local/lib -lBulletDynamics -lBulletCollision -lLinearMath -lBullet3Common
#include "bullet.h"
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/quat"
)

type Compound struct {
	shapeBase
	object.Component

	shapes []Shape
}

func NewCompound() *Compound {
	cmp := object.NewComponent(&Compound{
		shapeBase: newShapeBase(CompoundShape),
	})

	cmp.handle = C.goNewCompoundShape((*C.char)(unsafe.Pointer(cmp)))

	runtime.SetFinalizer(cmp, func(c *Compound) {
		c.destroy()
	})

	return cmp
}

func (c *Compound) OnEnable() {
	c.shapes = object.GetAllInChildren[Shape](c)
	for _, shape := range c.shapes {
		pos := c.Transform().Unproject(shape.Transform().WorldPosition())
		rot := quat.Ident()
		C.goAddChildShape(c.handle, shape.shape(), vec3ptr(&pos), quatPtr(&rot))
	}
	c.OnChange().Emit(c)
}

func (c *Compound) destroy() {
	c.shapes = nil
	if c.handle != nil {
		C.goDeleteShape(c.handle)
		c.handle = nil
	}
}
