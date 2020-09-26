package engine

// Application holds references to the basic engine components
type Application struct {
	Window   *Window
	Pipeline *Renderer
	Update   UpdateCallback
	Draw     RenderCallback
}

// NewApplication instantiates a new engine application
func NewApplication(title string, width, height int) *Application {
	highDpiEnabled := false
	wnd := CreateWindow(title, width, height, highDpiEnabled)

	// figure out render resolution if we're on a high dpi screen
	scale := float32(1.0)
	if highDpiEnabled {
		scale = wnd.Scale()
	}
	renderWidth, renderHeight := int32(float32(width)*scale), int32(float32(height)*scale)

	// set upp renderer
	renderer := NewRenderer()
	geoPass := NewGeometryPass(renderWidth, renderHeight)
	lightPass := NewLightPass(geoPass.Buffer)
	colorPass := NewColorPass(lightPass.Output, "saturated")
	renderer.Append("geometry", geoPass)
	renderer.Append("light", lightPass)
	renderer.Append("postprocess", colorPass)
	renderer.Append("output", NewOutputPass(colorPass.Output, geoPass.Buffer.Depth))
	renderer.Append("lines", NewLinePass())

	app := &Application{
		Window:   wnd,
		Pipeline: renderer,
	}

	// update callback
	app.Window.SetUpdateCallback(func(dt float32) {
		if app.Update != nil {
			app.Update(dt)
		}
	})

	// draw callback
	app.Window.SetRenderCallback(func(wnd *Window, dt float32) {
		if app.Draw != nil {
			app.Draw(wnd, dt)
		}
	})

	return app
}

// Run the application. Hands over control to the main window, looping until it is closed.
func (app *Application) Run() {
	app.Window.Loop()
}
