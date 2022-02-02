package shader

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

type ShaderID uint32

type T interface {
	Use()
	SetFragmentData(string)
	Attach(Stage)
	Link()
	Uniform(string) (UniformDesc, error)
	Attribute(string) (AttributeDesc, error)
	Mat4(string, mat4.T) error
	Vec2(string, vec2.T) error
	Vec3(string, vec3.T) error
	Vec4(string, vec4.T) error
	Vec3Array(string, []vec3.T) error
	Float(string, float32) error
	Int32(string, int) error
	Uint32(string, int) error
	Bool(string, bool) error
	RGB(string, color.T) error
	RGBA(string, color.T) error
	Texture2D(string, texture.Slot) error

	VertexPointers(interface{}) Pointers
}

type UniformMap map[string]UniformDesc
type AttributeMap map[string]AttributeDesc

type UniformDesc struct {
	Name  string
	Index int
	Size  int
	Type  types.Type
}

type AttributeDesc struct {
	Name  string
	Index int
	Size  int
	Type  types.Type
}
