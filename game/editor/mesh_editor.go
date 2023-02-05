package editor

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/material"
)

type meshEditor struct {
	T
}

func NewMeshEditor(ctx *Context, mesh mesh.T) *meshEditor {
	return object.New(&meshEditor{})
}

func init() {
	Register(mesh.New(mesh.Deferred, &material.Def{}), NewMeshEditor)
}
