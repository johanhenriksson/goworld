package ui

import (
    "github.com/johanhenriksson/goworld/render"
)

type Text struct {
    *Element
    Text    string
    Color   render.Color
}

func (m *Manager) NewText(text string, color render.Color, x, y, z float32) *Text {
    /* TODO: calculate size of text */
    w, h := float32(0.0), float32(0.0)
    el := m.NewElement(x,y,w,h,z)
    t := &Text {
        Element: el,
        Text: text,
        Color: color,
    }
    return t
}
