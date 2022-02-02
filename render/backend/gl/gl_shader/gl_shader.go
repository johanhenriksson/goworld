package gl_shader

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/johanhenriksson/goworld/render/backend/gl"
)

// Shader represents a GLSL program composed of several shaders
type glshader struct {
	ID    shader.ShaderID
	Name  string
	Debug bool

	shaders    []shader.Stage
	linked     bool
	uniforms   shader.UniformMap
	attributes shader.AttributeMap
}

// New creates a new shader program
func New(name string) shader.T {
	id := gl.CreateProgram()
	return &glshader{
		ID:         id,
		Name:       name,
		linked:     false,
		shaders:    make([]shader.Stage, 0),
		uniforms:   make(shader.UniformMap),
		attributes: make(shader.AttributeMap),
	}
}

// Compile compiles a set of GLSL files into a linked shader program.
// Filenames ending in vs, fs, gs indicate vertex, fragment and geometry shaders.
func Compile(name string, fileNames ...string) shader.T {
	shader := New(name)
	for _, fileName := range fileNames {
		stage := StageFromFile(fileName)
		shader.Attach(stage)
	}
	shader.Link()
	return shader
}

// Use binds the program for use in rendering
func (s *glshader) Use() {
	if !s.linked {
		panic(fmt.Sprintf("shader %s is not linked", s.Name))
	}
	gl.UseProgram(s.ID)
	if s.Debug {
		fmt.Println("use shader", s.Name)
	}
}

// SetFragmentData sets the name of the fragment color output variable
func (s *glshader) SetFragmentData(fragVariable string) {
	if err := gl.BindFragDataLocation(s.ID, fragVariable); err != nil {
		panic(err)
	}
}

// Attach a shader to the program. Panics if the program is already linked
func (s *glshader) Attach(stage shader.Stage) {
	if s.linked {
		panic(fmt.Sprintf("cant attach, shader %s is already linked", s.Name))
	}
	gl.AttachShader(s.ID, stage.ID())
	s.shaders = append(s.shaders, stage)
}

// Link the currently attached shaders into a program. Panics on failure
func (s *glshader) Link() {
	if s.linked {
		return
	}

	if err := gl.LinkProgram(s.ID); err != nil {
		panic(fmt.Sprintf("failed to compile %s: %s", s.Name, err))
	}

	s.readAttributes()
	s.readUniforms()

	s.linked = true
}

// Uniform returns a GLSL uniform location, and a bool indicating whether it exists or not
func (s *glshader) Uniform(uniform string) (shader.UniformDesc, error) {
	desc, ok := s.uniforms[uniform]
	if !ok {
		return shader.UniformDesc{
			Name:  uniform,
			Index: -1,
		}, fmt.Errorf("%w: %s", shader.ErrUnknownUniform, uniform)
	}
	return desc, nil
}

// Attribute returns a GLSL attribute location, and a bool indicating whether it exists or not
func (s *glshader) Attribute(attr string) (shader.AttributeDesc, error) {
	desc, ok := s.attributes[attr]
	if !ok {
		return shader.AttributeDesc{
			Name:  attr,
			Index: -1,
		}, fmt.Errorf("%w: %s", shader.ErrUnknownAttribute, attr)
	}
	return desc, nil
}

// Mat4 Sets a 4 by 4 matrix uniform value
func (s *glshader) Mat4(name string, mat4 mat4.T) error {
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformMatrix4f(s.ID, uniform, mat4)
	} else {
		return err
	}
}

// Vec2 sets a Vec2 uniform value
func (s *glshader) Vec2(name string, vec vec2.T) error {
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec2f(s.ID, uniform, vec)
	} else {
		return err
	}
}

// Vec3 sets a Vec3 uniform value
func (s *glshader) Vec3(name string, vec vec3.T) error {
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec3f(s.ID, uniform, vec)
	} else {
		return err
	}
}

func (s *glshader) Vec3Array(name string, vecs []vec3.T) error {
	if uniform, err := s.Uniform(fmt.Sprintf("%s[0]", name)); err == nil {
		return gl.UniformVec3fArray(s.ID, uniform, vecs)
	} else {
		return err
	}
}

// Vec4 sets a Vec4f uniform value
func (shader *glshader) Vec4(name string, vec vec4.T) error {
	if uniform, err := shader.Uniform(name); err == nil {
		return gl.UniformVec4f(shader.ID, uniform, vec)
	} else {
		return err
	}
}

// Int32 sets an integer 32 uniform value
func (s *glshader) Int32(name string, value int) error {
	if uniform, err := s.Uniform(name); err == nil {
		if name == "depth" {
			fmt.Println(uniform)
		}
		return gl.UniformVec1i(s.ID, uniform, value)
	} else {
		return err
	}
}

// Uint32 sets an unsigned integer 32 uniform value
func (s *glshader) Uint32(name string, value int) error {
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec1ui(s.ID, uniform, value)
	} else {
		return err
	}
}

func (s *glshader) Bool(name string, value bool) error {
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformBool(s.ID, uniform, value)
	} else {
		return err
	}
}

// Float sets a float uniform value
func (s *glshader) Float(name string, value float32) error {
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec1f(s.ID, uniform, value)
	} else {
		return err
	}
}

// RGB sets a uniform to a color RGB value
func (s *glshader) RGB(name string, color color.T) error {
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec3f(s.ID, uniform, vec3.T{
			X: color.R,
			Y: color.G,
			Z: color.B,
		})
	} else {
		return err
	}
}

// RGBA sets a uniform to a color RGBA value
func (s *glshader) RGBA(name string, color color.T) error {
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec4f(s.ID, uniform, vec4.T{
			X: color.R,
			Y: color.G,
			Z: color.B,
			W: color.A,
		})
	} else {
		return err
	}
}

func (s *glshader) Texture2D(name string, slot texture.Slot) error {
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformTexture2D(s.ID, uniform, slot)
	} else {
		return err
	}
}

func (s *glshader) VertexPointers(data interface{}) shader.Pointers {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Slice {
		return nil
	}

	el := t.Elem()
	if el.Kind() != reflect.Struct {
		fmt.Println("not a struct")
		return nil
	}

	size := int(el.Size())

	offset := 0
	pointers := make(shader.Pointers, 0, el.NumField())
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		tagstr := f.Tag.Get("vtx")
		if tagstr == "skip" {
			continue
		}
		tag, err := vertex.ParseTag(tagstr)
		if err != nil {
			fmt.Printf("tag error on %s.%s: %s\n", el.String(), f.Name, err)
			continue
		}

		gltype, err := gl.TypeFromString(tag.Type)
		if err != nil {
			panic(fmt.Errorf("invalid GL type: %s", tag.Type))
		}

		attr, err := s.Attribute(tag.Name)
		if errors.Is(err, shader.ErrUnknownAttribute) {
			// fmt.Printf("attribute %s does not exist on %s\n", tag.Name, shader.Name)
			offset += gltype.Size() * tag.Count
			continue
		}

		ptr := Pointer{
			Index:       int(attr.Index),
			Name:        tag.Name,
			Source:      gltype,
			Destination: gl.TypeCast(attr.Type),
			Elements:    tag.Count,
			Normalize:   tag.Normalize,
			Offset:      offset,
			Stride:      size,
		}

		pointers = append(pointers, ptr)

		offset += gltype.Size() * tag.Count
	}

	return pointers
}

func (s *glshader) readAttributes() {
	count := gl.GetActiveAttributeCount(s.ID)

	s.attributes = make(shader.AttributeMap, count)
	for i := 0; i < count; i++ {
		attr := s.readAttribute(i)
		s.attributes[attr.Name] = attr
		fmt.Println(s.Name, "uniform", attr.Name, attr.Type, "=", attr.Index, "size:", attr.Size)
	}
}

func (s *glshader) readAttribute(index int) shader.AttributeDesc {
	return gl.GetActiveAttribute(s.ID, index)
}

func (s *glshader) readUniforms() {
	count := gl.GetActiveUniformCount(s.ID)
	s.uniforms = make(map[string]shader.UniformDesc, count)
	for i := 0; i < count; i++ {
		uniform := s.readUniform(i)
		s.uniforms[uniform.Name] = uniform
		fmt.Println(s.Name, "uniform", uniform.Name, uniform.Type, "=", uniform.Index, "size:", uniform.Size)
	}
}

func (s *glshader) readUniform(index int) shader.UniformDesc {
	return gl.GetActiveUniform(s.ID, index)
}
