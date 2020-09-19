package engine

// Application holds references to the basic engine components
type Application struct {
	Window     *Window
	Scene      *Scene
	Render     *Renderer
	UpdateFunc UpdateCallback
}

// NewApplication instantiates a new engine application
func NewApplication(title string, width, height int) *Application {
	highDpiEnabled := false
	wnd := CreateWindow(title, width, height, highDpiEnabled)

	// create a scene
	scene := NewScene()

	// figure out render resolution if we're on a high dpi screen
	scale := float32(1.0)
	if highDpiEnabled {
		scale = wnd.Scale()
	}
	renderWidth, renderHeight := int32(float32(width)*scale), int32(float32(height)*scale)

	// set upp renderer
	renderer := NewRenderer(renderWidth, renderHeight, scene)
	geoPass := NewGeometryPass(renderWidth, renderHeight)
	lightPass := NewLightPass(geoPass.Buffer)
	colorPass := NewColorPass(lightPass.Output, "none")
	renderer.Append("geometry", geoPass)
	renderer.Append("light", lightPass)
	renderer.Append("postprocess", colorPass)
	renderer.Append("output", NewOutputPass(colorPass.Output))
	renderer.Append("lines", NewLinePass())

	app := &Application{
		Window: wnd,
		Scene:  scene,
		Render: renderer,
	}

	// update callback
	app.Window.SetUpdateCallback(func(dt float32) {
		app.Render.Update(dt)
		if app.UpdateFunc != nil {
			app.UpdateFunc(dt)
		}
	})

	// draw callback
	app.Window.SetRenderCallback(func(wnd *Window, dt float32) {
		app.Render.Draw()
	})

	return app
}

// Run the application. Hands over control to the main window, looping until it is closed.
func (app *Application) Run() {
	app.Window.Loop()
}
