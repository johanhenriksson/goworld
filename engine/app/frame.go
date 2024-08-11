package app

import (
	osimage "image"
	"runtime"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render/image"
)

// Render a single frame and return it as *image.RGBA
func Frame(args Args, scenefuncs ...object.SceneFunc) *osimage.RGBA {
	runtime.LockOSThread()
	args.Defaults()

	app := engine.New("goworld", 0)
	defer app.Destroy()

	buffer := engine.NewColorTarget(app.Device(), "output", image.FormatRGBA8Unorm, engine.TargetSize{
		Width:  args.Width,
		Height: args.Height,
		Frames: 1,
		Scale:  1,
	})
	defer buffer.Destroy()

	// create renderer
	renderer := args.Renderer(app, buffer)
	defer renderer.Destroy()

	// create scene
	scene := object.Scene(scenefuncs...)

	scene.Update(scene, 0)
	renderer.Draw(scene, 0, 0)

	return renderer.Screengrab()
}
