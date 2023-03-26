package propedit

import (
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/gui/node"
)

func Transform(key string, tf transform.T) node.T {
	return Container(key, []node.T{
		Vec3Field("position", "Position", Vec3Props{
			Value:    tf.Position(),
			OnChange: tf.SetPosition,
		}),
		Vec3Field("rotation", "Rotation", Vec3Props{
			Value:    tf.Rotation(),
			OnChange: tf.SetRotation,
		}),
		Vec3Field("scale", "Scale", Vec3Props{
			Value:    tf.Scale(),
			OnChange: tf.SetScale,
		}),
	})
}
