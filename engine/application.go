package engine

import (
    "github.com/johanhenriksson/goworld/ui"
)

type Application struct {
    Window      *Window
    Scene       *Scene
    Render      *Renderer
    UI          *ui.Manager
}

func NewApplication(title string, width, height int) *Application {
    wnd := CreateWindow(title, width, height)

    /* Scene */
    scene := NewScene()

    /* Renderer */
    renderer := NewRenderer(int32(width), int32(height), scene)

    /* UI Manager */
    uimgr := ui.NewManager(float32(width), float32(height))

    app := &Application {
        Window: wnd,
        Scene: scene,
        Render: renderer,
        UI: uimgr,
    }

    /* Update callback */
    app.Window.SetUpdateCallback(app.Render.Update)


    return app
}

func (app *Application) Run() {
    app.Window.Loop()
}
