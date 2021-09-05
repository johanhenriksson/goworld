package keys

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Action glfw.Action

const (
	Press   Action = Action(glfw.Press)
	Release        = Action(glfw.Release)
	Repeat         = Action(glfw.Repeat)
	Char           = Action(3)
)

func (a Action) String() string {
	switch a {
	case Press:
		return "Press"
	case Release:
		return "Release"
	case Char:
		return "Character"
	}
	return "Invalid"
}
