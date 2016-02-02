package render

import (
    "fmt"
    "strings"

	"github.com/go-gl/gl/v4.1-core/gl"
    mgl "github.com/go-gl/mathgl/mgl32"

    "github.com/johanhenriksson/goworld/util"
)

const (
    UnknownAttribute    AttributeLocation = -1
    UnknownUniform      UniformLocation   = -1
)

type AttributeLocation  int32
type UniformLocation    int32
type UniformMap         map[string]UniformLocation
type AttributeMap       map[string]AttributeLocation

/* Represents a GLSL program composed of several shaders. 
   Use CreateProgram() to instantiate */
type ShaderProgram struct {
    Id          uint32
    shaders     []*Shader
    linked      bool
    uniforms    UniformMap
    attributes  AttributeMap
}

func CreateProgram() *ShaderProgram {
    id := gl.CreateProgram()
    return &ShaderProgram {
        Id:         id,
        linked:     false,
        shaders:    make([]*Shader, 0),
        uniforms:   make(UniformMap),
        attributes: make(AttributeMap),
    }
}

/* Shorthand to compile a vertex & fragment shader and link them into a shader program.
   Uses the given file path plus ".vs.glsl" for the vertex shader and ".fs.glsl" for
   the fragment shader. */
func CompileVFShader(shaderFileName string) *ShaderProgram {
    program := CreateProgram()
    program.Attach(VertexShader(fmt.Sprintf("%s.vs.glsl", shaderFileName)))
    program.Attach(FragmentShader(fmt.Sprintf("%s.fs.glsl", shaderFileName)))
    program.Link()
    return program
}

/* Binds the program for use in rendering */
func (program *ShaderProgram) Use() {
    if !program.linked {
        panic("Shader program is not yet linked")
    }
	gl.UseProgram(program.Id)
}

/* Sets the name of the fragment color output variable */
func (program *ShaderProgram) SetFragmentData(fragVariable string) {
	gl.BindFragDataLocation(program.Id, 0, util.GLString(fragVariable))
}

/* Attach a shader to the program. Panics if the program is already linked */
func (program *ShaderProgram) Attach(shader *Shader) {
    if program.linked {
        panic("Cannot attach shader, program is already linked")
    }
    gl.AttachShader(program.Id, shader.Id)
    program.shaders = append(program.shaders, shader)
}

func (program *ShaderProgram) Link() {
    if program.linked {
        return
    }

	gl.LinkProgram(program.Id)

    /* Read status */
	var status int32
	gl.GetProgramiv(program.Id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program.Id, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program.Id, logLength, nil, gl.Str(log))

		panic(fmt.Sprintf("Failed to link program: %v", log))
	}

    program.linked = true
}

/* Returns a GLSL uniform location. If it doesnt exist, UnknownUniform is returned */
func (program *ShaderProgram) GetUniformLoc(uniform string) UniformLocation {
    loc, ok := program.uniforms[uniform]
    if !ok {
        loc = UniformLocation(gl.GetUniformLocation(program.Id, util.GLString(uniform)))
        if loc == UnknownUniform {
            panic("Unknown uniform: " + uniform)
        }
        program.uniforms[uniform] = loc
    }
    return loc
}

/* Returns a GLSL attribute location. If it doesnt exist, UnknownAttribute is returned */
func (program *ShaderProgram) GetAttrLoc(attr string) AttributeLocation {
    loc, ok := program.attributes[attr]
    if !ok {
        loc = AttributeLocation(gl.GetAttribLocation(program.Id, util.GLString(attr)))
        if loc == UnknownAttribute {
            panic("Unknown attribute: " + attr)
        }
        program.attributes[attr] = loc
    }
    return loc
}

/* Sets a 3x3 matrix uniform */
func (program *ShaderProgram) Matrix3f(name string, ptr *float32) {
    loc := program.GetUniformLoc(name)
    gl.UniformMatrix3fv(int32(loc), 1, false, ptr)
}

/* Sets a 4 by 4 matrix uniform */
func (program *ShaderProgram) Matrix4f(name string, ptr *float32) {
    loc := program.GetUniformLoc(name)
    gl.UniformMatrix4fv(int32(loc), 1, false, ptr)
}

/* Sets a Vec2 uniform */
func (program *ShaderProgram) Vec2(name string, vec *mgl.Vec2) {
    loc := program.GetUniformLoc(name)
    gl.Uniform2f(int32(loc), vec[0], vec[1])
}

/* Sets a Vec3 uniform */
func (program *ShaderProgram) Vec3(name string, vec *mgl.Vec3) {
    loc := program.GetUniformLoc(name)
    gl.Uniform3f(int32(loc), vec[0], vec[1], vec[2])
}

/* Sets a Vec4 uniform */
func (program *ShaderProgram) Vec4(name string, vec *mgl.Vec4) {
    loc := program.GetUniformLoc(name)
    gl.Uniform4f(int32(loc), vec[0], vec[1], vec[2], vec[3])
}

/* Sets an integer 32 uniform */
func (program *ShaderProgram) Int32(name string, val int32) {
    loc := program.GetUniformLoc(name)
    gl.Uniform1i(int32(loc), val)
}

/* Sets an unsigned integer 32 uniform */
func (program *ShaderProgram) UInt32(name string, val uint32) {
    loc := program.GetUniformLoc(name)
    gl.Uniform1ui(int32(loc), val)
}

/* Sets a float uniform */
func (program *ShaderProgram) Float(name string, val float32) {
    loc := program.GetUniformLoc(name)
    gl.Uniform1f(int32(loc), val)
}