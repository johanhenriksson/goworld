package engine

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/render"
)

type SceneFunc func(Renderer, object.T)
type RendererFunc func() Renderer

type Args struct {
	Title    string
	Width    int
	Height   int
	Backend  window.GlfwBackend
	Renderer RendererFunc
}

func Run(args Args, scenefuncs ...SceneFunc) {
	log.Println("goworld")

	go RunProfilingServer(6060)
	interrupt := NewInterrupter()

	// default to opengl backend
	backend := args.Backend
	if backend == nil {
		// default to opengl backend
		backend = &window.OpenGLBackend{}
	}
	defer backend.Destroy()

	// default to deferred opengl renderer
	if args.Renderer == nil {
		args.Renderer = NewRenderer
	}

	// create a window
	wnd, err := window.New(backend, window.Args{
		Title:  args.Title,
		Width:  args.Width,
		Height: args.Height,
	})
	if err != nil {
		panic(err)
	}

	var renderer Renderer
	recreateRenderer := func() {
		if renderer != nil {
			renderer.Destroy()
		}
		renderer = args.Renderer()
	}
	recreateRenderer()
	defer func() {
		renderer.Destroy()
	}()

	// create scene
	scene := object.New("Scene")
	wnd.SetInputHandler(scene)
	for _, scenefunc := range scenefuncs {
		scenefunc(renderer, scene)
	}

	// run the render loop
	log.Println("ready")

	for interrupt.Running() && !wnd.ShouldClose() {
		wnd.Poll()

		w, h := wnd.Size()
		screen := render.Screen{
			Width:  w,
			Height: h,
			Scale:  wnd.Scale(),
		}

		// update scene
		scene.Update(0.016)

		// find the first active camera
		camera := query.New[camera.T]().First(scene)
		if camera == nil {
			fmt.Println("no active camera in the scene")
			continue
		}

		// draw
		context, err := backend.Aquire()
		if err != nil {
			log.Println("swapchain recreated?? recreating renderer")
			recreateRenderer()
			continue
		}

		args := CreateRenderArgs(screen, camera)
		args.Context = context

		renderer.Draw(args, scene)
		backend.Present()
	}
}
