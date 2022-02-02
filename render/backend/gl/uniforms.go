package gl

import (
	"fmt"
	"strings"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func GetActiveUniform(id shader.ShaderID, index int) shader.UniformDesc {
	var gltype uint32
	var length, size int32
	buffer := strings.Repeat("\x00", 64)
	bufferPtr := gl.Str(buffer)
	gl.GetActiveUniform(uint32(id), uint32(index), int32(len(buffer))-1, &length, &size, &gltype, bufferPtr)
	loc := gl.GetUniformLocation(uint32(id), bufferPtr)

	return shader.UniformDesc{
		Name:  buffer[:length],
		Index: int(loc),
		Size:  int(size),
		Type:  Type(gltype).Cast(),
	}
}

func GetActiveUniformCount(id shader.ShaderID) int {
	var uniforms int32
	gl.GetProgramiv(uint32(id), gl.ACTIVE_UNIFORMS, &uniforms)
	return int(uniforms)
}

func uniformTypeError(uniform shader.UniformDesc, actual types.Type) error {
	return fmt.Errorf("%w: cant assign %s to uniform %s, expected %s",
		shader.ErrUniformType,
		actual,
		uniform.Name,
		uniform.Type)
}

func uniformUpdateErrorCheck() error {
	if err := gl.GetError(); err != None {
		return fmt.Errorf("%w: %d", shader.ErrUpdateUniform, err)
	}
	return nil
}

func UniformMatrix4f(id shader.ShaderID, uniform shader.UniformDesc, mat4 mat4.T) error {
	if uniform.Type != types.Mat4f {
		return uniformTypeError(uniform, types.Mat4f)
	}
	gl.ProgramUniformMatrix4fv(uint32(id), int32(uniform.Index), 1, false, &mat4[0])
	return uniformUpdateErrorCheck()
}

func UniformVec1f(id shader.ShaderID, uniform shader.UniformDesc, value float32) error {
	if uniform.Type != types.Float {
		return uniformTypeError(uniform, types.Float)
	}
	gl.ProgramUniform1f(uint32(id), int32(uniform.Index), value)
	return uniformUpdateErrorCheck()
}

func UniformVec2f(id shader.ShaderID, uniform shader.UniformDesc, vec vec2.T) error {
	if uniform.Type != types.Vec2f {
		return uniformTypeError(uniform, types.Vec2f)
	}
	gl.ProgramUniform2f(uint32(id), int32(uniform.Index), vec.X, vec.Y)
	return uniformUpdateErrorCheck()
}

func UniformVec3f(id shader.ShaderID, uniform shader.UniformDesc, vec vec3.T) error {
	if uniform.Type != types.Vec3f {
		return uniformTypeError(uniform, types.Vec3f)
	}
	gl.ProgramUniform3f(uint32(id), int32(uniform.Index), vec.X, vec.Y, vec.Z)
	return uniformUpdateErrorCheck()
}

func UniformVec3fArray(id shader.ShaderID, uniform shader.UniformDesc, vecs []vec3.T) error {
	if uniform.Type != types.Vec3f {
		return uniformTypeError(uniform, types.Vec3f)
	}
	if uniform.Size == 1 {
		return fmt.Errorf("%w: %s is not an array", shader.ErrUniformType, uniform.Name)
	}
	if len(vecs) >= uniform.Size {
		return fmt.Errorf("%w: array is too long for %s, max length: %d", shader.ErrUniformType, uniform.Name, uniform.Size)
	}
	for i, vec := range vecs {
		gl.ProgramUniform3f(uint32(id), int32(uniform.Index+i), vec.X, vec.Y, vec.Z)
	}
	return uniformUpdateErrorCheck()
}

func UniformVec4f(id shader.ShaderID, uniform shader.UniformDesc, vec vec4.T) error {
	if uniform.Type != types.Vec4f {
		return uniformTypeError(uniform, types.Vec4f)
	}
	gl.ProgramUniform4f(uint32(id), int32(uniform.Index), vec.X, vec.Y, vec.Z, vec.W)
	return uniformUpdateErrorCheck()
}

func UniformVec1ui(id shader.ShaderID, uniform shader.UniformDesc, value int) error {
	if uniform.Type != types.UInt32 {
		return uniformTypeError(uniform, types.UInt32)
	}
	gl.ProgramUniform1ui(uint32(id), int32(uniform.Index), uint32(value))
	return uniformUpdateErrorCheck()
}

func UniformVec1i(id shader.ShaderID, uniform shader.UniformDesc, value int) error {
	if uniform.Type != types.Int32 {
		return uniformTypeError(uniform, types.Int32)
	}
	gl.ProgramUniform1i(uint32(id), int32(uniform.Index), int32(value))
	return uniformUpdateErrorCheck()
}

func UniformBool(id shader.ShaderID, uniform shader.UniformDesc, value bool) error {
	if uniform.Type != types.Bool {
		return uniformTypeError(uniform, types.Bool)
	}
	iv := int32(0)
	if value {
		iv = 1
	}
	gl.ProgramUniform1i(uint32(id), int32(uniform.Index), iv)
	return uniformUpdateErrorCheck()
}

func UniformTexture2D(id shader.ShaderID, uniform shader.UniformDesc, slot texture.Slot) error {
	if uniform.Type != types.Texture2D {
		return uniformTypeError(uniform, types.Bool)
	}
	gl.ProgramUniform1i(uint32(id), int32(uniform.Index), int32(slot))
	return uniformUpdateErrorCheck()
}
