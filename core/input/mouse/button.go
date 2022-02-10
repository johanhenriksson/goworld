package mouse

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Button glfw.MouseButton

const (
	Button1 Button = Button(glfw.MouseButton1)
	Button2        = Button(glfw.MouseButton2)
	Button3        = Button(glfw.MouseButton3)
	Button4        = Button(glfw.MouseButton4)
	Button5        = Button(glfw.MouseButton5)
)

func (b Button) String() string {
	return fmt.Sprintf("Button %d", int(b)+1)
}
