package engine

import (
    "github.com/johanhenriksson/goworld/ui"
)

type Application struct {
    Window      *Window
    Scene       *Scene
    Render      *Renderer
    UI          *ui.Manager
    UpdateFunc  UpdateCallback
}

func NewApplication(title string, width, height int) *Application {
    wnd := CreateWindow(title, width, height)

    /* Scene */
    scene := NewScene()

    /* Renderer */
    renderer := NewRenderer(int32(width), int32(height), scene)
    geom_pass := NewGeometryPass(int32(width), int32(width))
    renderer.Append("geometry", geom_pass)
    renderer.Append("light", NewLightPass(geom_pass.Buffer))
    renderer.Append("lines", NewLinePass())

    /* UI Manager */
    uimgr := ui.NewManager(float32(width), float32(height))

    app := &Application {
        Window: wnd,
        Scene: scene,
        Render: renderer,
        UI: uimgr,
    }

    /* Update callback */
    app.Window.SetUpdateCallback(func(dt float32) {
        app.Render.Update(dt)
        if app.UpdateFunc != nil {
            app.UpdateFunc(dt)
        }
        inputEndFrame()
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
