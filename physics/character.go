package physics

import (
	"log"
	"runtime"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Character struct {
	object.Component

	handle   characterHandle
	shape    Shape
	step     float32
	world    *World
	grounded bool
}

func NewCharacter(height, radius, stepHeight float32) *Character {
	shape := NewCapsule(height, radius)
	handle := character_new(shape.shape(), stepHeight)
	character := object.NewComponent(&Character{
		handle: handle,
		shape:  shape,
		step:   stepHeight,
	})
	runtime.SetFinalizer(character, func(c *Character) {
		character_delete(&c.handle)
	})
	return character
}

func (c *Character) pullState() {
	// pull physics state
	state := character_state_pull(c.handle)

	c.Transform().SetPosition(state.position)
	c.Transform().SetRotation(state.rotation)
	c.grounded = state.grounded
}

func (c *Character) pushState() {
	// todo: not required unless we changed something
	// todo: include movement dir?
	character_state_push(c.handle, c.Transform().WorldPosition(), c.Transform().WorldRotation())
}

// Move the character controller. Called every frame to apply movement.
func (c *Character) Move(dir vec3.T) {
	dir.Scale(0.016)
	character_move(c.handle, dir)
}

// Jump applies a jumping force to the character
// todo: configurable?
func (c *Character) Jump() {
	character_jump(c.handle)
}

// Grounded returns true if the character is in contact with ground.
func (c *Character) Grounded() bool {
	return c.grounded
}

func (c *Character) OnEnable() {
	if c.world = object.GetInParents[*World](c); c.world != nil {
		c.world.AddCharacter(c)
	} else {
		log.Println("Character", c.Name(), ": No physics world in parents")
	}
}

func (c *Character) OnDisable() {
	if c.world != nil {
		c.world.RemoveCharacter(c)
		c.world = nil
	}
}
