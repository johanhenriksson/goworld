package hooks

import (
	"github.com/johanhenriksson/goworld/core/scene"
)

var sceneRef scene.T = nil

func SetScene(scene scene.T) {
	sceneRef = scene
}

func UseScene() scene.T {
	return sceneRef
}
