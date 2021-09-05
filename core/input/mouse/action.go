package mouse

import "github.com/go-gl/glfw/v3.1/glfw"

type Action int

const (
	Press   Action = Action(glfw.Press)
	Release        = Action(glfw.Release)
	Move           = Action(4)
	Scroll         = Action(5)
)

func (a Action) String() string {
	switch a {
	case Press:
		return "Press"
	case Release:
		return "Release"
	case Move:
		return "Move"
	case Scroll:
		return "Scroll"
	}
	return "Invalid"
}
