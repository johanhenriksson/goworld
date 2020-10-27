package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/ui"
)

func DebugBufferWindows(app *engine.Application) ui.Component {
	lightPass := app.Pipeline.Light
	geoPass := app.Pipeline.Geometry
	bufferWindows := ui.NewRect(ui.Style{"spacing": ui.Float(10)},
		newBufferWindow("Diffuse", geoPass.Buffer.Diffuse, false),
		newBufferWindow("Normal", geoPass.Buffer.Normal, false),
		newBufferWindow("Position", geoPass.Buffer.Position, false),
		newBufferWindow("Occlusion", lightPass.SSAO.Gaussian.Output, true),
		newBufferWindow("Light", lightPass.Output.Texture, false))
	bufferWindows.SetPosition(vec2.New(10, 10))
	bufferWindows.Flow(vec2.New(500, 1000))
	return bufferWindows
}

func newBufferWindow(title string, texture *render.Texture, depth bool) ui.Component {
	var img *ui.Image
	size := vec2.New(240, 160)
	if depth {
		img = ui.NewDepthImage(texture, size, true)
	} else {
		img = ui.NewImage(texture, size, true, ui.NoStyle)
	}

	return ui.NewRect(WindowStyle,
		ui.NewText(title, ui.NoStyle),
		img)
}
