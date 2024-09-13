package assets

import (
	"github.com/johanhenriksson/goworld/assets/fs"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/shader"
)

type Shader interface {
	Asset

	LoadShader(fs.Filesystem, *device.Device) *shader.Shader
}

var _ Shader = shader.Ref("")
