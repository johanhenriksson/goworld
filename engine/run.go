package engine

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/renderer"
	"github.com/johanhenriksson/goworld/geometry/gizmo/mover"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type SceneFunc func(renderer.T, object.T)
type RendererFunc func(vulkan.Target) renderer.T

type PreDrawable interface {
	object.Component
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
	debug.SetGCPercent(-1)

	go RunProfilingServer(6060)
	interrupt := NewInterrupter()

	backend := vulkan.New("goworld", 0)
	defer backend.Destroy()

	if args.Renderer == nil {
		panic("no renderer given")
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
	scene := object.New("Scene")

	// create editor
	editor := object.New("Editor")
	editor.Attach(collider.NewManager())
	mv := mover.New(mover.Args{})
	mv.Transform().SetPosition(vec3.New(1, 10, 1))
	editor.Adopt(mv)
	scene.Adopt(editor)

	// create game scene root
	game := object.New("Game")
	scene.Adopt(game)

	wnd.SetInputHandler(scene)

	for _, scenefunc := range scenefuncs {
		scenefunc(renderer, game)
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
		camera := query.New[camera.T]().First(scene)
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
		objects := query.New[PreDrawable]().Collect(scene)
		for _, component := range objects {
			component.PreDraw(args.Apply(component.Object().Transform().World()), scene)
		}

		// draw
		renderer.Draw(args, scene)

		// wait for submissions
		// wnd.Transferer().Wait()

		// present image
		wnd.Worker(context.Index).
			Present(wnd.Swapchain(), context)

		wnd.Worker(context.Index).Wait()

		// gc pass
		// this might be a decent place to run GC?
		// or a horrible one since we are waiting for vulkan stuff to complete
		collectGarbage()

		timing := counter.Sample()
		log.Printf(
			"frame: %2.fms, avg: %.2fms, peak: %.2f, fps: %.1f\n",
			1000*timing.Current,
			1000*timing.Average,
			1000*timing.Max,
			1.0/timing.Average)
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
		log.Printf("slow GC cycle: ran gc cycle in %.2fms", elapsed.Seconds()*1000)
	}
}
