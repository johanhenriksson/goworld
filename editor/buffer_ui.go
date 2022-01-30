package editor

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/ui"
)

func DebugBufferWindows(renderer *engine.Renderer) ui.Component {
	light := renderer.Light
	geom := renderer.Geometry
	ssao := renderer.SSAO
	bufferWindows := ui.NewRect(ui.Style{"spacing": ui.Float(10)},
		newBufferWindow("Diffuse", geom.Buffer.Diffuse(), false),
		newBufferWindow("Normal", geom.Buffer.Normal(), false),
		newBufferWindow("Position", geom.Buffer.Position(), false),
		newBufferWindow("Shadow", light.Shadows.Output, true),
		newBufferWindow("SSAO", ssao.Gaussian.Output, true))
	bufferWindows.SetPosition(vec2.New(10, 10))
	bufferWindows.Flow(vec2.New(500, 1000))
	return bufferWindows
}

func newBufferWindow(title string, texture texture.T, depth bool) ui.Component {
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
