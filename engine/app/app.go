package app

import (
	"log"
	"runtime"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/window"
	"github.com/johanhenriksson/goworld/engine/window/glfw"
)

func Run(args Args, scenefuncs ...object.SceneFunc) {
	runtime.LockOSThread()
	args.Defaults()

	go engine.RunProfilingServer(6060)
	interrupt := NewInterrupter()

	app := engine.New("goworld", 0)
	defer app.Destroy()

	// create a window
	wnd, err := glfw.NewWindow(app.Instance(), app.Device(), window.WindowArgs{
		Title:  args.Title,
		Width:  args.Width,
		Height: args.Height,
		Frames: 3,
	})
	if err != nil {
		panic(err)
	}
	defer wnd.Destroy()

	// create renderer
	renderer := args.Renderer(app, wnd)
	defer renderer.Destroy()

	// create scene
	scene := object.Scene(scenefuncs...)
	wnd.SetInputHandler(scene)

	object.Attach(scene, engine.NewStatsGUI())

	// run the render loop
	log.Println("ready")

	counter := engine.NewFrameCounter(60)
	for interrupt.Running() && !wnd.ShouldClose() {
		// update scene
		wnd.Poll()
		counter.Update()
		scene.Update(scene, counter.Delta())

		// draw
		renderer.Draw(scene, counter.Elapsed(), counter.Delta())
	}
}
