package render

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type DrawPass int32

const (
	GeometryPass DrawPass = iota
	LightPass
	LinePass
)

/** UI Component render interface */
type Drawable interface {
	Draw(DrawArgs)

	ZIndex() float32

	/* Render tree */
	Parent() Drawable
	SetParent(Drawable)
	Children() []Drawable
}

/** Passed to Drawables on render */
type DrawArgs struct {
	VP         mgl.Mat4
	MVP        mgl.Mat4
	Projection mgl.Mat4
	View       mgl.Mat4
	Transform  mgl.Mat4
	Pass       DrawPass
}
