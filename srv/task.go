package srv

import "log"

type Task interface {
	Start(Actor)
	Step(Actor, float32) bool

	// Stop is called by the controller if the task is interrupted
	Stop(Actor)
}

type TaskQueue struct {
	task Task
	todo []Task
}

func (c *TaskQueue) Queue(t Task) {
	c.todo = append(c.todo, t)
}

func (c *TaskQueue) Start(u Actor) {}
func (c *TaskQueue) Stop(u Actor)  {}

func (c *TaskQueue) Step(u Actor, dt float32) bool {
	if c.task == nil {
		if len(c.todo) > 0 {
			c.task = c.todo[0]
			c.todo = c.todo[1:]
			log.Println("unit", u.Name(), "starting task", c.task)
			c.task.Start(u)
		} else {
			log.Println("unit", u.Name(), "is idle")
			return true
		}
	}

	log.Println("unit", u.Name(), "tick", dt)

	if c.task.Step(u, dt) {
		// done
		c.task = nil
	}

	return len(c.todo) == 0
}

type TaskLoop struct {
	todo  []Task
	index int
}

func NewTaskLoop(tasks ...Task) *TaskLoop {
	return &TaskLoop{todo: tasks}
}

func (c *TaskLoop) Start(u Actor) {
	c.index = 0
}

func (c *TaskLoop) Stop(u Actor) {}

func (c *TaskLoop) Step(u Actor, dt float32) bool {
	if len(c.todo) == 0 {
		return true
	}

	if c.todo[c.index].Step(u, dt) {
		c.index = (c.index + 1) % len(c.todo)
	}

	return false
}
