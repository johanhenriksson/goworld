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

type Args struct {
	Title    string
	Width    int
	Height   int
	Backend  vulkan.T
	Renderer RendererFunc
}

func Run(args Args, scenefuncs ...SceneFunc) {
	log.Println("goworld")

	// disable automatic garbage collection
	debug.SetGCPercent(-1)

	go RunProfilingServer(6060)
	interrupt := NewInterrupter()

	backend := args.Backend
	if backend == nil {
		panic("no backend provided")
	}
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

	var renderer renderer.T
	recreateRenderer := func() {
		if renderer != nil {
			renderer.Destroy()
		}
		renderer = args.Renderer(wnd)
	}
	recreateRenderer()
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

	lastFrameTime := time.Now()
	framesSinceGC := 0
	for interrupt.Running() && !wnd.ShouldClose() {
		wnd.Poll()

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

		// draw
		context, err := wnd.Aquire()
		if err != nil {
			log.Println("swapchain recreated?? recreating renderer")
			recreateRenderer()
			continue
		}

		args := createRenderArgs(screen, camera)
		args.Context = context

		renderer.Draw(args, scene)
		wnd.Present()

		// update scene
		endFrameTime := time.Now()
		elapsed := endFrameTime.Sub(lastFrameTime)
		scene.Update(float32(elapsed.Seconds()))

		remainingTime := float32(1.0/60 - time.Since(lastFrameTime).Seconds())
		lastFrameTime = endFrameTime

		if remainingTime > 0.001 || framesSinceGC > 60 {
			// manually trigger garbage collection
			// log.Printf("garbage collection pass r=%.2fms f=%d\n", 1000*remainingTime, framesSinceGC)
			runtime.GC()
			framesSinceGC = 0
		} else {
			framesSinceGC++
		}
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
