package game

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

func CreateScene(scene object.T) (*Player, *Chunk) {
	scene.Attach(light.NewDirectional(light.DirectionalArgs{
		Intensity: 1.6,
		Color:     color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
		Direction: vec3.New(0.95, -1.9, 1.05),
		Shadows:   true,
	}))

	// create chunk
	world := NewWorld(31481234, 16)
	chunk := world.AddChunk(0, 0)

	// first person controls
	player := NewPlayer(vec3.New(-5, 5, -5), func(player *Player, target vec3.T) (bool, vec3.T) {
		height := world.HeightAt(target)
		if target.Y < height {
			return true, vec3.New(target.X, height, target.Z)
		}
		return false, vec3.Zero
	})
	player.Flying = true
	player.Eye.Transform().SetRotation(vec3.New(0, 0, 0))
	scene.Adopt(player)

	return player, chunk
}
