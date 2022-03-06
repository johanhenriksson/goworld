package engine

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/render"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

type SceneFunc func(Renderer, object.T)
type RendererFunc func() Renderer

type Args struct {
	Title     string
	Width     int
	Height    int
	SceneFunc SceneFunc
	Backend   window.GlfwBackend
	Renderer  RendererFunc
}

func Run(args Args) {
	fmt.Println("goworld")

	// cpu profiling
	flag.Parse()
	if *cpuprofile != "" {
		os.MkdirAll("profiling", 0755)
		ppath := fmt.Sprintf("profiling/%s", *cpuprofile)
		f, err := os.Create(ppath)
		if err != nil {
			panic(err)
		}
		fmt.Println("writing cpu profiling output to", ppath)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	running := true
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		for range sigint {
			if !running {
				log.Println("Kill")
				os.Exit(1)
			} else {
				log.Println("Interrupt")
				running = false
			}
		}
	}()

	backend := args.Backend
	if backend == nil {
		// default to opengl backend
		backend = &window.OpenGLBackend{}
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
	if args.Renderer != nil {
		renderer = args.Renderer()
	} else {
		// default to deferred opengl renderer
		renderer = NewRenderer()
	}
	defer renderer.Destroy()

	// create scene
	scene := object.New("Scene")
	wnd.SetInputHandler(scene)
	args.SceneFunc(renderer, scene)

	// run the render loop
	fmt.Println("ready")

	for running && !wnd.ShouldClose() {
		wnd.Poll()

		w, h := wnd.Size()
		screen := render.Screen{
			Width:  w,
			Height: h,
			Scale:  wnd.Scale(),
		}

		// update scene
		scene.Update(0.030)

		// find the first active camera
		camera := query.New[camera.T]().First(scene)
		if camera == nil {
			fmt.Println("no active camera in the scene")
			continue
		}

		// draw
		context, err := backend.Aquire()
		if err != nil {
			continue
		}

		args := CreateRenderArgs(screen, camera)
		args.Context = context

		renderer.Draw(args, scene)
		backend.Present()
	}
}
