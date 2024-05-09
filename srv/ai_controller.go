package srv

import (
	"log"
	"time"

	"github.com/johanhenriksson/goworld/math/vec3"
)

type Actor interface {
	Observable

	ID() Identity
	Name() string
	Area() Area

	Position() vec3.T
	SetPosition(p vec3.T)

	Target() Identity
	SetTarget(t Identity)

	Spawn(Area, Identity, vec3.T)
	Despawn()
}

type Behavior map[string]Task

type AIController struct {
	target   Actor
	task     Task
	behavior Behavior
}

func NewAIController(target Actor, tasks Behavior) *AIController {
	aic := &AIController{
		target:   target,
		behavior: tasks,
	}
	go aic.loop()
	return aic
}

func (c *AIController) loop() {
	// subscribe to the target unit events
	events := EventBuffer()
	unsub := c.target.Subscribe(c, func(ev Event) {
		events <- ev
	})
	defer unsub()

	// set up task ticker
	lastTick := time.Now()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		// react to events from the target unit
		case ev := <-events:
			switch ev := ev.(type) {
			case PositionUpdateEvent:
				// react to the unit position update
				// move towards the unit
				log.Println("unit", c.target.Name(), "position update", ev.Position)

			case BehaviorEvent:
				// react to behavior events
				log.Println("unit", c.target.Name(), "behavior event", ev.Behavior)
				c.Do(ev.Behavior)

			default:
				log.Println("unit", c.target.Name(), "unhandled event", ev)
			}

		case now := <-ticker.C:
			// execute ai logic
			dt := now.Sub(lastTick)
			lastTick = now

			if c.task == nil {
				c.task = c.behavior["idle"]
				if c.task == nil {
					log.Println("unit", c.target.Name(), "has no idle task")
					return
				}
			}

			log.Println("unit", c.target.Name(), "tick", dt)

			if c.task.Step(c.target, float32(dt.Seconds())) {
				// done
				c.task = nil
			}
		}
	}
}

func (c *AIController) Do(behavior string) {
	if task, ok := c.behavior[behavior]; ok {
		c.task = task
	}
}

func (c *AIController) Stop() {
	if c.task != nil {
		c.task.Stop(c.target)
	}
	c.task = nil
}

type BehaviorEvent struct {
	source   any
	Behavior string
}

func (e BehaviorEvent) Source() any {
	return e.source
}
