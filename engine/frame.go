package engine

import (
	osimage "image"
	"runtime"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/graph"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

// Render a single frame and return it as *image.RGBA
func Frame(args Args, scenefuncs ...SceneFunc) *osimage.RGBA {
	runtime.LockOSThread()

	backend := vulkan.New("goworld", 0)
	defer backend.Destroy()

	if args.Renderer == nil {
		args.Renderer = graph.Default
	}

	buffer := vulkan.NewColorTarget(backend.Device(), "output", image.FormatRGBA8Unorm, vulkan.TargetSize{
		Width:  args.Width,
		Height: args.Height,
		Frames: 1,
		Scale:  1,
	})
	defer buffer.Destroy()

	// create renderer
	renderer := args.Renderer(backend, buffer)
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
