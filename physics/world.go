package physics

import (
    "math"
    "github.com/ianremmler/ode"
)

type World struct {
    gravity  float32
    world    ode.World
    space    ode.Space
    contacts ode.JointGroup
    timestep float64
}

func init() {
    /* Initialize Open Dynamics Engine */
    ode.Init(0, ode.AllAFlag)
}

func NewWorld() *World {
    world := ode.NewWorld()
    space := ode.NilSpace().NewHashSpace()

    w := &World {
        gravity: -9.82,
        timestep: 0.01,
        world: world,
        space: space,
        contacts: ode.NewJointGroup(10000),
    }

    world.SetGravity(ode.V3(0, float64(w.gravity), 0))
    world.SetCFM(1.0e-5)
    world.SetERP(0.2)
    world.SetContactSurfaceLayer(0.001)
    world.SetContactMaxCorrectingVelocity(0.9)
    world.SetAutoDisable(true)

    return w
}

func (w *World) Update() {
    w.space.Collide(0, w.nearCallback)
    w.world.QuickStep(w.timestep)
    w.contacts.Empty()
}

func (w *World) nearCallback(data interface{}, obj1, obj2 ode.Geom) {
    body1 := obj1.Body()
    body2 := obj2.Body()

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
}

func (w *World) NewPlane(x,y,z,c float32) {
    w.space.NewPlane(ode.V4(float64(x), float64(y), float64(z), float64(c)))
}
