package propedit

import (
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func Transform(key string, tf transform.T) node.T {
	return Container(key, []node.T{
		Vec3Field("position", "Position", Vec3Props{
			Value:    tf.Position(),
			OnChange: tf.SetPosition,
		}),
		Vec3Field("rotation", "Rotation", Vec3Props{
			Value: tf.Rotation().Euler(),
			OnChange: func(euler vec3.T) {
				tf.SetRotation(quat.Euler(euler.X, euler.Y, euler.Z))
			},
		}),
		Vec3Field("scale", "Scale", Vec3Props{
			Value:    tf.Scale(),
			OnChange: tf.SetScale,
		}),
	})
}
