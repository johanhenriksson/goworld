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

	// set upp renderer
	// this belongs somewhere else probably
	// actually, the entire application concept is pretty much rendundant at this point.
	// perhaps the window should be passed directly to a renderer?
	// or the other way around?
	renderer := NewRenderer()

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

	app.Window.SetResizeCallback(func(wnd *Window, width, height int) {
		renderer.Resize(width, height)
	})

	return app
}

// Run the application. Hands over control to the main window, looping until it is closed.
func (app *Application) Run() {
	app.Window.Loop()
}
