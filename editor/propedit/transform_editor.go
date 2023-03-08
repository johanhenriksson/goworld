package propedit

import (
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/gui/node"
)

func Transform(key string, tf transform.T) node.T {
	return Container(key, []node.T{
		Field("position", "Position", []node.T{
			Vec3("position", Vec3Props{
				Value:    tf.Position(),
				OnChange: tf.SetPosition,
			}),
		}),
		Field("rotation", "Rotation", []node.T{
			Vec3("rotation", Vec3Props{
				Value:    tf.Rotation(),
				OnChange: tf.SetRotation,
			}),
		}),
		Field("scale", "Scale", []node.T{
			Vec3("scale", Vec3Props{
				Value:    tf.Scale(),
				OnChange: tf.SetScale,
			}),
		}),
	})
}
