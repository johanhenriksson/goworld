package physics

import (
	"fmt"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/ode"
	"math"
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

	w := &World{
		gravity:  -9.82,
		timestep: 0.01,
		world:    world,
		space:    space,
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
		/* contact info */
		contact := ode.NewContact()
		contact.Surface.Mode = ode.BounceCtParam | ode.SoftCFMCtParam
		contact.Surface.Mu = math.Inf(1)
		contact.Surface.Mu2 = 0
		contact.Surface.Bounce = 0.01
		contact.Surface.BounceVel = 0.1
		contact.Geom = ct

		/* add contact joint until next frame */
		ctj := w.world.NewContactJoint(w.contacts, contact)
		ctj.Attach(body1, body2)

		/* collision events */
		event_contact := Contact{
			Position: FromOdeVec3(ct.Pos),
			Normal:   FromOdeVec3(ct.Normal),
			Depth:    float32(ct.Depth),
		}

		if obj1.Data() != nil && obj1.Data() != nil {
			col1 := obj1.Data().(Collider)
			col2 := obj2.Data().(Collider)
			col1.OnCollision(col2, event_contact)
			col2.OnCollision(col1, event_contact)
		}
	}
}

type RaycastHit struct {
	Position mgl.Vec3
	Normal   mgl.Vec3
	Distance float32
	Collider Collider
}

func (w *World) Raycast(length float32, origin, direction mgl.Vec3) []RaycastHit {
	hits := make([]RaycastHit, 0, 2)

	ray := w.space.NewRay(float64(length))
	ray.SetPosDir(ToOdeVec3(origin), ToOdeVec3(direction))
	ray.SetData("ray")

	w.space.Collide(0, func(data interface{}, obj1, obj2 ode.Geom) {
		var other Collider

		if obj1 == ray {
			other = obj2.Data().(Collider)
		} else if obj2 == ray {
			other = obj1.Data().(Collider)
		} else {
			return
		}

		cts := obj1.Collide(obj2, 1, 0)
		if len(cts) > 0 {
			ct := cts[0]

			/* collision events */
			hit := RaycastHit{
				Position: FromOdeVec3(ct.Pos),
				Normal:   FromOdeVec3(ct.Normal),
				Distance: float32(ct.Depth),
				Collider: other,
			}
			hits = append(hits, hit)

			fmt.Println("hit", other, fmt.Sprintf("at [%.1f,%.1f,%.1f] normal: [%.1f,%.1f,%.1f], depth: %.1f",
				hit.Position[0], hit.Position[1], hit.Position[2],
				hit.Normal[0], hit.Normal[1], hit.Normal[2],
				hit.Distance))
		}
	})
	ray.Destroy()
	return hits
}
