package render

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/util"
)

// TODO: return proper errors, dont just crash

type ShaderInput struct {
	Name  string
	Index int32
	Type  GLType
}

// ShaderInputs maps names to GL attribute &uniform locations
type ShaderInputs map[string]ShaderInput

// Shader represents a GLSL program composed of several shaders
type Shader struct {
	ID    uint32
	Name  string
	Debug bool

	shaders    []*ShaderStage
	linked     bool
	uniforms   ShaderInputs
	attributes ShaderInputs
}

// CreateShader creates a new shader program
func CreateShader(name string) *Shader {
	id := gl.CreateProgram()
	return &Shader{
		ID:         id,
		Name:       name,
		linked:     false,
		shaders:    make([]*ShaderStage, 0),
		uniforms:   make(ShaderInputs),
		attributes: make(ShaderInputs),
	}
}

// CompileShader compiles a set of GLSL files into a linked shader program.
// Filenames ending in vs, fs, gs indicate vertex, fragment and geometry shaders.
func CompileShader(name string, fileNames ...string) *Shader {
	shader := CreateShader(name)
	for _, fileName := range fileNames {
		stage := CompileShaderStage(fileName)
		shader.Attach(stage)
	}
	shader.Link()
	return shader
}

// Use binds the program for use in rendering
func (program *Shader) Use() {
	if !program.linked {
		panic(fmt.Sprintf("shader %s is not linked", program.Name))
	}
	gl.UseProgram(program.ID)
}

// SetFragmentData sets the name of the fragment color output variable
func (shader *Shader) SetFragmentData(fragVariable string) {
	cstr, free := util.GLString(fragVariable)
	gl.BindFragDataLocation(shader.ID, 0, *cstr)
	free()
}

// Attach a shader to the program. Panics if the program is already linked
func (shader *Shader) Attach(stage *ShaderStage) {
	if shader.linked {
		panic(fmt.Sprintf("cant attach, shader %s is already linked", shader.Name))
	}
	gl.AttachShader(shader.ID, stage.ID)
	shader.shaders = append(shader.shaders, stage)
}

// Link the currently attached shaders into a program. Panics on failure
func (shader *Shader) Link() {
	if shader.linked {
		return
	}

	gl.LinkProgram(shader.ID)

	/* Read status */
	var status int32
	gl.GetProgramiv(shader.ID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shader.ID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shader.ID, logLength, nil, gl.Str(log))

		panic(fmt.Sprintf("failed to link program %s: %v", shader.Name, log))
	}

	shader.readAttributes()
	shader.readUniforms()

	shader.linked = true
}

// Uniform returns a GLSL uniform location, and a bool indicating whether it exists or not
func (shader *Shader) Uniform(uniform string) (ShaderInput, bool) {
	input, ok := shader.uniforms[uniform]
	if !ok {
		if shader.Debug {
			fmt.Println("Unknown uniform", uniform, "in shader", shader.Name)
		}
		return ShaderInput{Name: uniform, Index: -1}, false
	}
	return input, true
}

// Attribute returns a GLSL attribute location, and a bool indicating whether it exists or not
func (shader *Shader) Attribute(attr string) (ShaderInput, bool) {
	input, ok := shader.attributes[attr]
	if !ok {
		if shader.Debug {
			fmt.Println("Unknown attribute", attr, "in shader", shader.Name)
		}
		return ShaderInput{Name: attr, Index: -1}, false
	}
	return input, true
}

// Mat4 Sets a 4 by 4 matrix uniform value
func (shader *Shader) Mat4(name string, mat4 *mat4.T) {
	if input, ok := shader.Uniform(name); ok {
		if input.Type != Mat4f {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Mat4f, name, input.Type))
		}
		gl.UniformMatrix4fv(input.Index, 1, false, &mat4[0])
	}
}

// Vec2 sets a Vec2 uniform value
func (shader *Shader) Vec2(name string, vec *vec2.T) {
	if input, ok := shader.Uniform(name); ok {
		if input.Type != Vec3f {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Vec2f, name, input.Type))
		}
		gl.Uniform2f(input.Index, vec.X, vec.Y)
	}
}

// Vec3 sets a Vec3 uniform value
func (program *Shader) Vec3(name string, vec *vec3.T) {
	if input, ok := program.Uniform(name); ok {
		if input.Type != Vec3f {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Vec3f, name, input.Type))
		}
		gl.Uniform3f(input.Index, vec.X, vec.Y, vec.Z)
	}
}

// Vec4 sets a Vec4f uniform value
func (shader *Shader) Vec4(name string, vec *vec4.T) {
	if input, ok := shader.Uniform(name); ok {
		if input.Type != Vec3f {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Vec4f, name, input.Type))
		}
		gl.Uniform4f(input.Index, vec.X, vec.Y, vec.Z, vec.W)
	}
}

// Int32 sets an integer 32 uniform value
func (shader *Shader) Int32(name string, val int32) {
	if input, ok := shader.Uniform(name); ok {
		if input.Type != Int32 && input.Type != Texture2D {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Int32, name, input.Type))
		}
		gl.Uniform1i(input.Index, val)
	}
}

// UInt32 sets an unsigned integer 32 uniform value
func (shader *Shader) UInt32(name string, val uint32) {
	if input, ok := shader.Uniform(name); ok {
		if input.Type != UInt32 {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", UInt32, name, input.Type))
		}
		gl.Uniform1ui(input.Index, val)
	}
}

// Float sets a float uniform value
func (shader *Shader) Float(name string, val float32) {
	if input, ok := shader.Uniform(name); ok {
		if input.Type != Float {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Float, name, input.Type))
		}
		gl.Uniform1f(input.Index, val)
	}
}

// RGB sets a uniform to a color RGB value
func (shader *Shader) RGB(name string, color Color) {
	if input, ok := shader.Uniform(name); ok {
		if input.Type != Vec3f {
			panic(fmt.Errorf("cant assign RGB color to uniform %s, expected %s", name, input.Type))
		}
		gl.Uniform3f(input.Index, color.R, color.G, color.B)
	}
}

// RGBA sets a uniform to a color RGBA value
func (shader *Shader) RGBA(name string, color Color) {
	if input, ok := shader.Uniform(name); ok {
		if input.Type != Vec4f {
			panic(fmt.Errorf("cant assign RGBA color to uniform %s, expected %s", name, input.Type))
		}
		gl.Uniform4f(input.Index, color.R, color.G, color.B, color.A)
	}
}

func (shader *Shader) readAttributes() {
	var attributes int32
	gl.GetProgramiv(shader.ID, gl.ACTIVE_ATTRIBUTES, &attributes)

	shader.attributes = make(ShaderInputs, int(attributes))
	for i := 0; i < int(attributes); i++ {
		attr := shader.readAttribute(i)
		shader.attributes[attr.Name] = attr
		fmt.Println(shader.Name, "attrib", attr.Name, attr.Type)
	}
}

func (shader *Shader) readAttribute(index int) ShaderInput {
	var gltype uint32
	var length, size int32
	buffer := strings.Repeat("\x00", 64)
	bufferPtr := gl.Str(buffer)
	gl.GetActiveAttrib(shader.ID, uint32(index), int32(len(buffer))-1, &length, &size, &gltype, bufferPtr)
	loc := gl.GetAttribLocation(shader.ID, bufferPtr)

	return ShaderInput{
		Name:  buffer[:length],
		Index: loc,
		Type:  GLType(gltype),
	}
}

func (shader *Shader) readUniforms() {
	var uniforms int32
	gl.GetProgramiv(shader.ID, gl.ACTIVE_UNIFORMS, &uniforms)

	shader.uniforms = make(ShaderInputs, int(uniforms))
	for i := 0; i < int(uniforms); i++ {
		uniform := shader.readUniform(i)
		shader.uniforms[uniform.Name] = uniform
		fmt.Println(shader.Name, "uniform", uniform.Name, uniform.Type)
	}
}

func (shader *Shader) readUniform(index int) ShaderInput {
	var gltype uint32
	var length, size int32
	buffer := strings.Repeat("\x00", 64)
	bufferPtr := gl.Str(buffer)
	gl.GetActiveUniform(shader.ID, uint32(index), int32(len(buffer))-1, &length, &size, &gltype, bufferPtr)
	loc := gl.GetUniformLocation(shader.ID, bufferPtr)

	return ShaderInput{
		Name:  buffer[:length],
		Index: loc,
		Type:  GLType(gltype),
	}
}
