package main

import (
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/editor"
	_ "github.com/johanhenriksson/goworld/editor/builtin"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/app"
	"github.com/johanhenriksson/goworld/game/chunk"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

func main() {
	app.Run(
		engine.Args{
			Width:  1600,
			Height: 1200,
			Title:  "voxel editor",
		},
		nil,
		editor.Scene(func(scene object.Object) {
			chk := chunk.New("editable", 16, 16, 16)

			object.Builder(object.Empty("Chunk")).
				Attach(chunk.NewMesh(chk)).
				Parent(scene).
				Create()

			object.Attach(scene,
				object.Builder(object.Empty("Sun")).
					Attach(light.NewDirectional(light.DirectionalArgs{
						Intensity: 1.3,
						Color:     color.RGB(1, 1, 1),
						Shadows:   true,
						Cascades:  4,
					})).
					Position(vec3.New(16, 16, -4)).
					Rotation(quat.Euler(40, 0, 0)).
					Create())
		}),
	)
}
