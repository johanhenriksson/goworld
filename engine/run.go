package engine

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type SceneFunc func(renderer.T, object.T)
type RendererFunc func(vulkan.Target) renderer.T

type PreDrawable interface {
	object.T
	PreDraw(render.Args, object.T) error
}

type Args struct {
	Title    string
	Width    int
	Height   int
	Renderer RendererFunc
}

func Run(args Args, scenefuncs ...SceneFunc) {
	log.Println("goworld")

	// disable automatic garbage collection
	// debug.SetGCPercent(-1)

	go RunProfilingServer(6060)
	interrupt := NewInterrupter()

	backend := vulkan.New("goworld", 0)
	defer backend.Destroy()

	if args.Renderer == nil {
		args.Renderer = renderer.NewGraph
	}

	// create a window
	wnd, err := backend.Window(vulkan.WindowArgs{
		Title:  args.Title,
		Width:  args.Width,
		Height: args.Height,
	})
	if err != nil {
		panic(err)
	}

	// create renderer
	renderer := args.Renderer(wnd)
	defer func() {
		renderer.Destroy()
	}()

	// create scene
	scene := object.Empty("Scene")
	wnd.SetInputHandler(scene)
	for _, scenefunc := range scenefuncs {
		scenefunc(renderer, scene)
	}

	// run the render loop
	log.Println("ready")

	counter := NewFrameCounter(60)

	currentTime := time.Now()
	// framesSinceGC := 0
	for interrupt.Running() && !wnd.ShouldClose() {
		newTime := time.Now()
		frameTime := newTime.Sub(currentTime)
		currentTime = newTime

		// update scene
		wnd.Poll()
		scene.Update(float32(frameTime.Seconds()))

		// render
		screen := render.Screen{
			Width:  wnd.Width(),
			Height: wnd.Height(),
			Scale:  wnd.Scale(),
		}

		// find the first active camera
		camera := object.Query[camera.T]().First(scene)
		if camera == nil {
			fmt.Println("no active camera in the scene")
			continue
		}

		// aquire next frame
		context, err := wnd.Aquire()
		if err != nil {
			log.Println("swapchain recreated?? recreating renderer")
			renderer.Recreate()
			continue
		}

		args := createRenderArgs(screen, camera)
		args.Context = context

		// pre-draw
		objects := object.Query[PreDrawable]().Collect(scene)
		for _, object := range objects {
			object.PreDraw(args.Apply(object.Transform().World()), scene)
		}

		// draw
		renderer.Draw(args, scene)

		// wait for submissions
		wnd.Transferer().Wait()

		// present image
		wnd.Worker(context.Index).
			Present(wnd.Swapchain(), context)

		wnd.Worker(context.Index).Wait()

		// gc pass
		// this might be a decent place to run GC?
		// or a horrible one since we are waiting for vulkan stuff to complete
		// collectGarbage()

		counter.Sample()
		// log.Printf(
		// 	"frame: %2.fms, avg: %.2fms, peak: %.2f, fps: %.1f\n",
		// 	1000*timing.Current,
		// 	1000*timing.Average,
		// 	1000*timing.Max,
		// 	1.0/timing.Average)
	}
}

func createRenderArgs(screen render.Screen, cam camera.T) render.Args {
	return render.Args{
		Projection: cam.Projection(),
		View:       cam.View(),
		VP:         cam.ViewProj(),
		MVP:        cam.ViewProj(),
		Transform:  mat4.Ident(),
		Position:   cam.Transform().WorldPosition(),
		Clear:      cam.ClearColor(),
		Forward:    cam.Transform().Forward(),
		Viewport:   screen,
	}
}

func collectGarbage() {
	start := time.Now()
	runtime.GC()
	elapsed := time.Since(start)
	if elapsed.Milliseconds() > 2 {
		log.Printf("slow GC cycle: %.2fms", elapsed.Seconds()*1000)
	}
}
