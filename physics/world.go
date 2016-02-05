package physics

import (
    "math"
    "github.com/ianremmler/ode"
)

type World struct {
    world    ode.World
    space    ode.Space
    contacts ode.JointGroup
}

func init() {
    /* Initialize Open Dynamics Engine */
    ode.Init(0, ode.AllAFlag)
}

func NewWorld() *World {
    world := ode.NewWorld()

    world.SetGravity(ode.V3(0,-1,0))
    world.SetCFM(1.0e-5)
    world.SetERP(0.2)
    world.SetContactSurfaceLayer(0.001)
    world.SetContactMaxCorrectingVelocity(0.9)
    world.SetAutoDisable(true)

    w := &World {
        world: world,
        space: ode.NilSpace().NewHashSpace(),
        contacts: ode.NewJointGroup(10000),
    }
    return w
}

func (w *World) Update() {
    w.space.Collide(0, func(data interface{}, obj1, obj2 ode.Geom) {
        body1, body2 := obj1.Body(), obj2.Body()
        cts := obj1.Collide(obj2, 1, 0)
        for _, ct := range cts {
            contact := ode.NewContact()
            contact.Surface.Mode = ode.BounceCtParam | ode.SoftCFMCtParam;
            contact.Surface.Mu = math.Inf(1)
            contact.Surface.Mu2 = 0;
            contact.Surface.Bounce = 0.01;
            contact.Surface.BounceVel = 0.1;

            contact.Geom = ct
            ct := w.world.NewContactJoint(w.contacts, contact)
            ct.Attach(body1, body2)
        }
    })
    w.world.QuickStep(0.01)
    w.contacts.Empty()

}

func (w *World) NewPlane(x,y,z,c float32) {
    w.space.NewPlane(ode.V4(float64(x), float64(y), float64(z), float64(c)))
}
