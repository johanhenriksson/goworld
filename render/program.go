package render

import (
    "fmt"
    "errors"
    "strings"

	"github.com/go-gl/gl/v4.1-core/gl"

    "github.com/johanhenriksson/goworld/util"
)

type AttributeLocation int32
type UniformLocation int32
type UniformMap map[string]UniformLocation
type AttributeMap map[string]AttributeLocation

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

func CompileVFShader(shaderFileName string) *ShaderProgram {
    program := CreateProgram()
    program.Attach(VertexShader(fmt.Sprintf("%s.vs.glsl", shaderFileName)))
    program.Attach(FragmentShader(fmt.Sprintf("%s.fs.glsl", shaderFileName)))
    if err := program.Link(); err != nil {
        panic(err)
    }
    return program
}

func (program *ShaderProgram) Use() {
    if !program.linked {
        panic("Shader program is not yet linked")
    }
	gl.UseProgram(program.Id)
}

func (program *ShaderProgram) SetFragmentData(fragVariable string) {
	gl.BindFragDataLocation(program.Id, 0, util.GLString(fragVariable))
}

func (program *ShaderProgram) Attach(shader *Shader) {
    gl.AttachShader(program.Id, shader.Id)
    program.shaders = append(program.shaders, shader)
}

func (program *ShaderProgram) Link() error {
	gl.LinkProgram(program.Id)

	var status int32
	gl.GetProgramiv(program.Id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program.Id, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program.Id, logLength, nil, gl.Str(log))

		return errors.New(fmt.Sprintf("Failed to link program: %v", log))
	}

    program.linked = true
    return nil
}

func (program *ShaderProgram) GetUniformLoc(uniform string) UniformLocation {
    loc, ok := program.uniforms[uniform]
    if !ok {
        loc = UniformLocation(gl.GetUniformLocation(program.Id, util.GLString(uniform)))
        if loc < 0 {
            panic("Uniform doesnt exist: " + uniform)
        }
        program.uniforms[uniform] = loc
    }
    return loc
}

func (program *ShaderProgram) GetAttrLoc(attr string) AttributeLocation {
    loc, ok := program.attributes[attr]
    if !ok {
        loc = AttributeLocation(gl.GetAttribLocation(program.Id, util.GLString(attr)))
        if loc < 0 {
            return -1
        }
        program.attributes[attr] = loc
    }
    return loc
}

func (program *ShaderProgram) Matrix4f(name string, ptr *float32) {
    loc := program.GetUniformLoc(name)
	gl.UniformMatrix4fv(int32(loc), 1, false, ptr)
}

func (program *ShaderProgram) Int32(name string, val int32) {
    loc := program.GetUniformLoc(name)
    gl.Uniform1i(int32(loc), val)
}

func (program *ShaderProgram) UInt32(name string, val uint32) {
    loc := program.GetUniformLoc(name)
    gl.Uniform1ui(int32(loc), val)
}
