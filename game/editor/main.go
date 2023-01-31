package editor

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/game"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/geometry/gizmo/mover"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

func Scene(render renderer.T, scene object.T) {
	// collision support
	object.Attach(scene, NewSelectManager())

	world := game.NewWorld(31481284, 8)
	chonk := world.AddChunk(0, 0)

	// first person controls
	player := game.NewPlayer(vec3.New(0, 20, -11), nil)
	player.Eye.Transform().SetRotation(vec3.New(-30, 0, 0))
	object.Attach(scene, player)

	voxedit := NewVoxelEditor(chonk, player.Camera, render)
	voxedit.SetActive(true)

	// chunk mesh & editor
	object.Builder(object.Empty("Chunk")).
		Attach(chunk.NewMesh(chonk)).
		Attach(voxedit).
		Attach(collider.NewBox(collider.Box{
			Center: vec3.New(4, 4, 4),
			Size:   vec3.New(8, 8, 8),
		})).
		Parent(scene).
		Create()

	// directional light
	object.Attach(scene, light.NewDirectional(light.DirectionalArgs{
		Intensity: 1.6,
		Color:     color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
		Direction: vec3.New(0.95, -1.9, 1.05),
		Shadows:   true,
	}))

	// mover gizmo
	object.Builder(mover.New(mover.Args{})).
		Position(vec3.New(1, 10, 1)).
		Parent(scene).
		Create()
}
