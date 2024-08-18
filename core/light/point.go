package light

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
)

type PointArgs struct {
	Color     color.T
	Range     float32
	Intensity float32
}

type Point struct {
	object.Component

	Color     object.Property[color.T]
	Range     object.Property[float32]
	Intensity object.Property[float32]
	Falloff   object.Property[float32]
}

var _ T = &Point{}

func init() {
	object.Register[*Point](object.TypeInfo{
		Name:        "Point Light",
		Deserialize: DeserializePoint,
		Create: func(pool object.Pool) (object.Component, error) {
			return NewPoint(pool, PointArgs{
				Color:     color.White,
				Range:     10,
				Intensity: 1,
			}), nil
		},
	})
}

func NewPoint(pool object.Pool, args PointArgs) *Point {
	return object.NewComponent(pool, &Point{
		Color:     object.NewProperty(args.Color),
		Range:     object.NewProperty(args.Range),
		Intensity: object.NewProperty(args.Intensity),
		Falloff:   object.NewProperty(float32(2)),
	})
}

func (lit *Point) Name() string      { return "PointLight" }
func (lit *Point) Type() Type        { return TypePoint }
func (lit *Point) CastShadows() bool { return false }

func (lit *Point) LightData(shadowmaps ShadowmapStore) uniform.Light {
	return uniform.Light{
		Type:      uint32(TypePoint),
		Position:  vec4.Extend(lit.Transform().WorldPosition(), 0),
		Color:     lit.Color.Get(),
		Intensity: lit.Intensity.Get(),
		Range:     lit.Range.Get(),
		Falloff:   lit.Falloff.Get(),
	}
}

func (lit *Point) Shadowmaps() int {
	return 0
}

func (lit *Point) ShadowProjection(mapIndex int) uniform.Camera {
	panic("todo")
}

type PointState struct {
	object.ComponentState
	PointArgs
}

func (lit *Point) Serialize(enc object.Encoder) error {
	if err := object.EncodeComponent(enc, lit.Component); err != nil {
		return err
	}
	return enc.Encode(*lit)
}

func DeserializePoint(pool object.Pool, dec object.Decoder) (object.Component, error) {
	cmp, err := object.DecodeComponent(pool, dec)
	if err != nil {
		return nil, err
	}

	lit := &Point{}
	if err := dec.Decode(lit); err != nil {
		return nil, err
	}
	lit.Component = cmp

	return object.NewComponent(pool, lit), nil
}
