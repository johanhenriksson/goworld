package render

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/johanhenriksson/goworld/util"
)

const (
	// UnknownAttribute is returned when an attribute name can't be resolved
	UnknownAttribute AttributeLocation = -1
	// UnknownUniform is returned when a uniform name can't be resolved
	UnknownUniform UniformLocation = -1
)

// TODO: return proper errors, dont just crash

// AttributeLocation is a GL attribute location index
type AttributeLocation int32

// UniformLocation is a GL uniform location index
type UniformLocation int32

// UniformMap maps uniform names to GL uniform locations
type UniformMap map[string]UniformLocation

// AttributeMap maps attribute names to GL attribute locations
type AttributeMap map[string]AttributeLocation

// ShaderProgram represents a GLSL program composed of several shaders
type ShaderProgram struct {
	ID         uint32
	shaders    []*Shader
	linked     bool
	uniforms   UniformMap
	attributes AttributeMap
}

// CreateProgram creates a new shader program
func CreateProgram() *ShaderProgram {
	id := gl.CreateProgram()
	return &ShaderProgram{
		ID:         id,
		linked:     false,
		shaders:    make([]*Shader, 0),
		uniforms:   make(UniformMap),
		attributes: make(AttributeMap),
	}
}

// CompileVFShader is a shorthand to compile a vertex & fragment shader and link them into a shader program.
// Uses the given file path plus ".vs.glsl" for the vertex shader and ".fs.glsl" for
// the fragment shader.
func CompileVFShader(shaderFileName string) *ShaderProgram {
	program := CreateProgram()
	program.Attach(VertexShader(fmt.Sprintf("%s.vs.glsl", shaderFileName)))
	program.Attach(FragmentShader(fmt.Sprintf("%s.fs.glsl", shaderFileName)))
	program.Link()
	return program
}

// Use binds the program for use in rendering
func (program *ShaderProgram) Use() {
	if !program.linked {
		panic("Shader program is not yet linked")
	}
	gl.UseProgram(program.ID)
}

// SetFragmentData sets the name of the fragment color output variable
func (program *ShaderProgram) SetFragmentData(fragVariable string) {
	cstr, free := util.GLString(fragVariable)
	gl.BindFragDataLocation(program.ID, 0, *cstr)
	free()
}

// Attach a shader to the program. Panics if the program is already linked
func (program *ShaderProgram) Attach(shader *Shader) {
	if program.linked {
		panic("Cannot attach shader, program is already linked")
	}
	gl.AttachShader(program.ID, shader.ID)
	program.shaders = append(program.shaders, shader)
}

// Link the currently attached shaders into a program. Panics on failure
func (program *ShaderProgram) Link() {
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

		panic(fmt.Sprintf("Failed to link program: %v", log))
	}

	program.linked = true
}

// GetUniformLoc returns a GLSL uniform location, and a bool indicating whether it exists or not
func (program *ShaderProgram) GetUniformLoc(uniform string) (UniformLocation, bool) {
	loc, ok := program.uniforms[uniform]
	if !ok {
		// get C string
		cstr, free := util.GLString(uniform)
		defer free()

		loc = UniformLocation(gl.GetUniformLocation(program.ID, *cstr))
		if loc == UnknownUniform {
			return loc, false
		}
		program.uniforms[uniform] = loc
	}
	return loc, true
}

// GetAttrLoc returns a GLSL attribute location, and a bool indicating whether it exists or not
func (program *ShaderProgram) GetAttrLoc(attr string) (AttributeLocation, bool) {
	loc, ok := program.attributes[attr]
	if !ok {
		// get c string
		cstr, free := util.GLString(attr)
		defer free()

		loc = AttributeLocation(gl.GetAttribLocation(program.ID, *cstr))
		if loc == UnknownAttribute {
			return loc, false
		}
		program.attributes[attr] = loc
	}
	return loc, true
}

// Mat3f sets a 3x3 matrix uniform value
func (program *ShaderProgram) Mat3f(name string, mat3 mgl.Mat3) {
	if loc, ok := program.GetUniformLoc(name); ok {
		gl.UniformMatrix3fv(int32(loc), 1, false, &mat3[0])
	}
}

// Mat4f Sets a 4 by 4 matrix uniform value
func (program *ShaderProgram) Mat4f(name string, mat4 mgl.Mat4) {
	if loc, ok := program.GetUniformLoc(name); ok {
		gl.UniformMatrix4fv(int32(loc), 1, false, &mat4[0])
	}
}

// Vec2 sets a Vec2 uniform value
func (program *ShaderProgram) Vec2(name string, vec *mgl.Vec2) {
	if loc, ok := program.GetUniformLoc(name); ok {
		gl.Uniform2f(int32(loc), vec[0], vec[1])
	}
}

// Vec3 sets a Vec3 uniform value
func (program *ShaderProgram) Vec3(name string, vec *mgl.Vec3) {
	if loc, ok := program.GetUniformLoc(name); ok {
		gl.Uniform3f(int32(loc), vec[0], vec[1], vec[2])
	}
}

// Vec4 sets a Vec4f uniform value
func (program *ShaderProgram) Vec4(name string, vec *mgl.Vec4) {
	if loc, ok := program.GetUniformLoc(name); ok {
		gl.Uniform4f(int32(loc), vec[0], vec[1], vec[2], vec[3])
	}
}

// Int32 sets an integer 32 uniform value
func (program *ShaderProgram) Int32(name string, val int32) {
	if loc, ok := program.GetUniformLoc(name); ok {
		gl.Uniform1i(int32(loc), val)
	}
}

// UInt32 sets an unsigned integer 32 uniform value
func (program *ShaderProgram) UInt32(name string, val uint32) {
	if loc, ok := program.GetUniformLoc(name); ok {
		gl.Uniform1ui(int32(loc), val)
	}
}

// Float sets a float uniform value
func (program *ShaderProgram) Float(name string, val float32) {
	if loc, ok := program.GetUniformLoc(name); ok {
		gl.Uniform1f(int32(loc), val)
	}
}

// RGB sets a uniform to a color RGB value
func (program *ShaderProgram) RGB(name string, color Color) {
	if loc, ok := program.GetUniformLoc(name); ok {
		gl.Uniform3f(int32(loc), color.R, color.G, color.B)
	}
}

// RGBA sets a uniform to a color RGBA value
func (program *ShaderProgram) RGBA(name string, color Color) {
	if loc, ok := program.GetUniformLoc(name); ok {
		gl.Uniform4f(int32(loc), color.R, color.G, color.B, color.A)
	}
}
