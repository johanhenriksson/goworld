package app

import (
	osimage "image"
	"runtime"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/graph"
	"github.com/johanhenriksson/goworld/render/image"
)

// Render a single frame and return it as *image.RGBA
func Frame(args engine.Args, graphFunc graph.GraphFunc, scenefuncs ...engine.SceneFunc) *osimage.RGBA {
	runtime.LockOSThread()

	app := engine.New("goworld", 0)
	defer app.Destroy()

	if graphFunc == nil {
		graphFunc = graph.Default
	}

	buffer := engine.NewColorTarget(app.Device(), "output", image.FormatRGBA8Unorm, engine.TargetSize{
		Width:  args.Width,
		Height: args.Height,
		Frames: 1,
		Scale:  1,
	})
	defer buffer.Destroy()

	// create renderer
	renderer := graphFunc(app, buffer)
	defer renderer.Destroy()

	// create scene
	scene := object.Empty("Scene")
	for _, scenefunc := range scenefuncs {
		scenefunc(scene)
	}

	scene.Update(scene, 0)
	renderer.Draw(scene, 0, 0)

	return renderer.Screengrab()
}
