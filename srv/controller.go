package srv

import (
	"log"
	"time"

	"github.com/johanhenriksson/goworld/math/vec3"
)

type Actor interface {
	Observable

	Name() string
	Position() vec3.T
	SetPosition(p vec3.T)
}

type AIController struct {
	task   Task
	todo   []Task
	target Actor
}

func NewAIController(target Actor, tasks []Task) *AIController {
	aic := &AIController{
		target: target,
		todo:   tasks,
	}
	go aic.loop()
	return aic
}

func (c *AIController) loop() {
	// subscribe to the target unit events
	events := make(chan Event, 1024)
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
			case UnitPositionUpdateEvent:
				// react to the unit position update
				// move towards the unit
				log.Println("unit", c.target.Name(), "position update", ev.Position)
			case TaskEvent:
				log.Println("unit", c.target.Name(), "aquired task", ev.Task)
				if c.task == nil {
					c.task = ev.Task
				} else {
					c.todo = append(c.todo, ev.Task)
				}

			default:
				log.Println("unit", c.target.Name(), "unhandled event", ev)
			}

		case now := <-ticker.C:
			// execute ai logic
			dt := now.Sub(lastTick)
			lastTick = now

			if c.task == nil {
				if len(c.todo) > 0 {
					c.task = c.todo[0]
					c.todo = c.todo[1:]
					log.Println("unit", c.target.Name(), "starting task", c.task)
					c.task.Start(c.target)
				} else {
					log.Println("unit", c.target.Name(), "is idle")
					continue
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
