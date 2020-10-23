package render

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/util"
)

const (
	// UnknownAttribute is returned when an attribute name can't be resolved
	UnknownAttribute int32 = -1
	// UnknownUniform is returned when a uniform name can't be resolved
	UnknownUniform int32 = -1
)

// TODO: return proper errors, dont just crash

// AttributeLocation is a GL attribute location index
type AttributeLocation int32

// UniformLocation is a GL uniform location index
type UniformLocation int32

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

// CompileShader is a shorthand to compile a vertex & fragment shader and link them into a shader program.
// Uses the given file path plus ".vs.glsl" for the vertex shader and ".fs.glsl" for
// the fragment shader.
func CompileShader(shaderFileName string) *Shader {
	name := filepath.Base(shaderFileName)
	program := CreateShader(name)
	program.Attach(VertexShader(fmt.Sprintf("%s.vs.glsl", shaderFileName)))
	program.Attach(FragmentShader(fmt.Sprintf("%s.fs.glsl", shaderFileName)))

	// optional geometry shader
	gsPath := fmt.Sprintf("%s.gs.glsl", shaderFileName)
	if _, err := os.Stat(gsPath); err == nil {
		program.Attach(GeometryShader(gsPath))
	}

	program.Link()

	program.getAttributes()
	program.getUniforms()

	return program
}

func CompileShaderFiles(name, path string, fileNames ...string) *Shader {
	program := CreateShader(name)
	for _, fileName := range fileNames {
		if len(fileName) < 3 {
			panic(fmt.Errorf("invalid shader filename: %s", fileName))
		}
		kind := fileName[len(fileName)-3:]
		switch kind {
		case ".fs":
			program.Attach(FragmentShader(fmt.Sprintf("%s/%s.glsl", path, fileName)))
		case ".vs":
			program.Attach(VertexShader(fmt.Sprintf("%s/%s.glsl", path, fileName)))
		case ".gs":
			program.Attach(GeometryShader(fmt.Sprintf("%s/%s.glsl", path, fileName)))
		default:
			panic(fmt.Errorf("invalid shader type %s: %s", kind, fileName))
		}
	}
	program.Link()

	program.getAttributes()
	program.getUniforms()

	return program
}

func (program *Shader) getAttributes() {
	var count int32
	gl.GetProgramiv(program.ID, gl.ACTIVE_ATTRIBUTES, &count)

	program.attributes = make(ShaderInputs)
	for i := 0; i < int(count); i++ {
		name, loc, gltype := program.readAttribute(i)
		program.attributes[name] = ShaderInput{
			Name:  name,
			Index: loc,
			Type:  gltype,
		}
		fmt.Printf("attribute %+v\n", program.attributes[name])
	}
}

func (program *Shader) readAttribute(index int) (string, int32, GLType) {
	var gltype uint32
	var length, size int32
	buffer := strings.Repeat("\x00", 129)
	bufferPtr := gl.Str(buffer)
	gl.GetActiveAttrib(program.ID, uint32(index), 128, &length, &size, &gltype, bufferPtr)
	loc := gl.GetAttribLocation(program.ID, bufferPtr)
	name := buffer[:length]
	return name, loc, GLType(gltype)
}

func (program *Shader) getUniforms() {
	var count int32
	gl.GetProgramiv(program.ID, gl.ACTIVE_UNIFORMS, &count)

	program.uniforms = make(ShaderInputs)
	for i := 0; i < int(count); i++ {
		name, loc, gltype := program.readUniform(i)
		program.uniforms[name] = ShaderInput{
			Name:  name,
			Index: loc,
			Type:  gltype,
		}
		fmt.Printf("uniform %+v\n", program.uniforms[name])
	}
}

func (program *Shader) readUniform(index int) (string, int32, GLType) {
	var gltype uint32
	var length, size int32
	buffer := strings.Repeat("\x00", 129)
	bufferPtr := gl.Str(buffer)
	gl.GetActiveUniform(program.ID, uint32(index), 128, &length, &size, &gltype, bufferPtr)
	loc := gl.GetUniformLocation(program.ID, bufferPtr)
	name := buffer[:length]
	return name, loc, GLType(gltype)
}

// Use binds the program for use in rendering
func (program *Shader) Use() {
	if !program.linked {
		panic(fmt.Sprintf("Shader program %s is not yet linked", program.Name))
	}
	gl.UseProgram(program.ID)
}

// SetFragmentData sets the name of the fragment color output variable
func (program *Shader) SetFragmentData(fragVariable string) {
	cstr, free := util.GLString(fragVariable)
	gl.BindFragDataLocation(program.ID, 0, *cstr)
	free()
}

// Attach a shader to the program. Panics if the program is already linked
func (program *Shader) Attach(shader *ShaderStage) {
	if program.linked {
		panic(fmt.Sprintf("Cannot attach to shader %s, program is already linked", program.Name))
	}
	gl.AttachShader(program.ID, shader.ID)
	program.shaders = append(program.shaders, shader)
}

// Link the currently attached shaders into a program. Panics on failure
func (program *Shader) Link() {
	if program.linked {
		return
	}

	gl.LinkProgram(program.ID)

	/* Read status */
	var status int32
	gl.GetProgramiv(program.ID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program.ID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program.ID, logLength, nil, gl.Str(log))

		panic(fmt.Sprintf("Failed to link program %s: %v", program.Name, log))
	}

	program.linked = true
}

// getUniform returns a GLSL uniform location, and a bool indicating whether it exists or not
func (program *Shader) getUniform(uniform string) (ShaderInput, bool) {
	input, ok := program.uniforms[uniform]
	if !ok {
		if program.Debug {
			fmt.Println("Unknown uniform", uniform, "in shader", program.Name)
		}
		return ShaderInput{Name: uniform, Index: -1}, false
	}
	return input, true
}

// getAttribute returns a GLSL attribute location, and a bool indicating whether it exists or not
func (program *Shader) getAttribute(attr string) (ShaderInput, bool) {
	input, ok := program.attributes[attr]
	if !ok {
		if program.Debug {
			fmt.Println("Unknown attribute", attr, "in shader", program.Name)
		}
		return ShaderInput{Name: attr, Index: -1}, false
	}
	return input, true
}

// Mat4 Sets a 4 by 4 matrix uniform value
func (program *Shader) Mat4(name string, mat4 *mat4.T) {
	if input, ok := program.getUniform(name); ok {
		if input.Type != Mat4f {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Mat4f, name, input.Type))
		}
		gl.UniformMatrix4fv(input.Index, 1, false, &mat4[0])
	}
}

// Vec2 sets a Vec2 uniform value
func (program *Shader) Vec2(name string, vec *vec2.T) {
	if input, ok := program.getUniform(name); ok {
		if input.Type != Vec3f {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Vec2f, name, input.Type))
		}
		gl.Uniform2f(input.Index, vec.X, vec.Y)
	}
}

// Vec3 sets a Vec3 uniform value
func (program *Shader) Vec3(name string, vec *vec3.T) {
	if input, ok := program.getUniform(name); ok {
		if input.Type != Vec3f {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Vec3f, name, input.Type))
		}
		gl.Uniform3f(input.Index, vec.X, vec.Y, vec.Z)
	}
}

// Vec4 sets a Vec4f uniform value
func (program *Shader) Vec4(name string, vec *vec4.T) {
	if input, ok := program.getUniform(name); ok {
		if input.Type != Vec3f {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Vec4f, name, input.Type))
		}
		gl.Uniform4f(input.Index, vec.X, vec.Y, vec.Z, vec.W)
	}
}

// Int32 sets an integer 32 uniform value
func (program *Shader) Int32(name string, val int32) {
	if input, ok := program.getUniform(name); ok {
		if input.Type != Int32 && input.Type != Texture2D {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Int32, name, input.Type))
		}
		gl.Uniform1i(input.Index, val)
	}
}

// UInt32 sets an unsigned integer 32 uniform value
func (program *Shader) UInt32(name string, val uint32) {
	if input, ok := program.getUniform(name); ok {
		if input.Type != UInt32 {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", UInt32, name, input.Type))
		}
		gl.Uniform1ui(input.Index, val)
	}
}

// Float sets a float uniform value
func (program *Shader) Float(name string, val float32) {
	if input, ok := program.getUniform(name); ok {
		if input.Type != Float {
			panic(fmt.Errorf("cant assign %s to uniform %s, expected %s", Float, name, input.Type))
		}
		gl.Uniform1f(input.Index, val)
	}
}

// RGB sets a uniform to a color RGB value
func (program *Shader) RGB(name string, color Color) {
	if input, ok := program.getUniform(name); ok {
		if input.Type != Vec3f {
			panic(fmt.Errorf("cant assign RGB color to uniform %s, expected %s", name, input.Type))
		}
		gl.Uniform3f(input.Index, color.R, color.G, color.B)
	}
}

// RGBA sets a uniform to a color RGBA value
func (program *Shader) RGBA(name string, color Color) {
	if input, ok := program.getUniform(name); ok {
		if input.Type != Vec4f {
			panic(fmt.Errorf("cant assign RGBA color to uniform %s, expected %s", name, input.Type))
		}
		gl.Uniform4f(input.Index, color.R, color.G, color.B, color.A)
	}
}
