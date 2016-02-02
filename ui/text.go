package ui

import (
    "github.com/johanhenriksson/goworld/render"
)

type Text struct {
    *Image
    Text    string
    Color   render.Color
}

func (m *Manager) NewText(text string, color render.Color, x, y, z float32) *Text {
    /* TODO: calculate size of text */
    w, h := float32(200.0), float32(25.0)
    fnt := render.LoadFont("assets/fonts/luximr.ttf", 16.0, 72.0, 1.5)
    texture := fnt.Render(text, w, h, color)
    img := m.NewImage(texture, x, y, w, h, z)

    return &Text {
        Image: img,
        Text: text,
        Color: color,
    }
}
