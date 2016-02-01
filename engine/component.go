package engine

import (
    "github.com/johanhenriksson/goworld/render"
)

type Component interface {
    Update(float32)
    Draw(render.DrawArgs)
}
