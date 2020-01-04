package engine

import (
	"github.com/johanhenriksson/goworld/ui"
)

type Application struct {
	Window     *Window
	Scene      *Scene
	Render     *Renderer
	UI         *ui.Manager
	UpdateFunc UpdateCallback
}

func NewApplication(title string, width, height int) *Application {
	highDpiEnabled := false
	wnd := CreateWindow(title, width, height, highDpiEnabled)

	/* Scene */
	scene := NewScene()

	// figure out render resolution
	scale := float32(1.0)
	if highDpiEnabled {
		scale = wnd.Scale()
	}
	renderWidth, renderHeight := int32(float32(width)*scale), int32(float32(height)*scale)

	/* Renderer */
	renderer := NewRenderer(renderWidth, renderHeight, scene)
	geoPass := NewGeometryPass(renderWidth, renderHeight)
	lightPass := NewLightPass(geoPass.Buffer)
	colorPass := NewColorPass(lightPass.Output, "saturated")
	renderer.Append("geometry", geoPass)
	renderer.Append("light", lightPass)
	renderer.Append("postprocess", colorPass)
	renderer.Append("output", NewOutputPass(colorPass.Output))
	//renderer.Append("lines", NewLinePass())

	/* UI Manager */
	uimgr := ui.NewManager(float32(width), float32(height))

	app := &Application{
		Window: wnd,
		Scene:  scene,
		Render: renderer,
		UI:     uimgr,
	}

	/* Update callback */
	app.Window.SetUpdateCallback(func(dt float32) {
		app.Render.Update(dt)
		if app.UpdateFunc != nil {
			app.UpdateFunc(dt)
		}
	})

	/* Draw callback */
	app.Window.SetRenderCallback(func(wnd *Window, dt float32) {
		/* render scene */
		app.Render.Draw()
		/* draw user interface */
		app.UI.Draw()
	})

	return app
}

/* Hands over control to the application. Will loop until the main window is closed */
func (app *Application) Run() {
	app.Window.Loop()
}
