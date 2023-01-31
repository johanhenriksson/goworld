package game

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

func CreateScene(render renderer.T, scene object.T) {
	object.Attach(scene, light.NewDirectional(light.DirectionalArgs{
		Intensity: 1.6,
		Color:     color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
		Direction: vec3.New(0.95, -1.9, 1.05),
		Shadows:   true,
	}))

	// create chunk
	world := NewWorld(31481284, 8)
	chonk := world.AddChunk(0, 0)

	// first person controls
	player := NewPlayer(vec3.New(0, 20, -11), func(player *Player, target vec3.T) (bool, vec3.T) {
		height := world.HeightAt(target)
		if target.Y < height {
			return true, vec3.New(target.X, height, target.Z)
		}
		return false, vec3.Zero
	})
	player.Eye.Transform().SetRotation(vec3.New(-30, 0, 0))
	object.Attach(scene, player)

	object.Builder(object.Empty("Chunk")).
		Attach(chunk.NewMesh(chonk)).
		Position(vec3.New(8, 0, 0)).
		// Attach(collider.NewBox(collider.Box{
		// 	Center: vec3.New(8, 8, 8),
		// 	Size:   vec3.New(16, 16, 16),
		// })).
		Parent(scene).
		Create()

	object.Builder(object.Empty("light1")).
		Position(vec3.New(10, 9, 13)).
		Attach(light.NewPoint(light.PointArgs{
			Attenuation: light.DefaultAttenuation,
			Color:       color.Red,
			Range:       15,
			Intensity:   15,
		})).
		Parent(scene).
		Create()

}
