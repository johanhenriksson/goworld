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

var activeShader shader.ShaderID = 0

// Shader represents a GLSL program composed of several shaders
type glshader struct {
	id         shader.ShaderID
	name       string
	shaders    []shader.Stage
	linked     bool
	uniforms   shader.UniformMap
	attributes shader.AttributeMap
	state      map[string]any
}

// New creates a new shader program
func New(name string) shader.T {
	id := gl.CreateProgram()
	if id == 0 {
		panic(fmt.Sprintf("failed to create shader %s", name))
	}
	return &glshader{
		id:         id,
		name:       name,
		linked:     false,
		shaders:    make([]shader.Stage, 0),
		uniforms:   make(shader.UniformMap),
		attributes: make(shader.AttributeMap),
		state:      make(map[string]any),
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
	if err := shader.Link(); err != nil {
		panic(err)
	}
	return shader
}

func (s *glshader) Name() string {
	return s.name
}

// Use binds the program for use in rendering
func (s *glshader) Use() error {
	if !s.linked {
		return fmt.Errorf("shader %s is not linked", s.name)
	}
	if activeShader != s.id {
		if err := gl.UseProgram(s.id); err != nil {
			return fmt.Errorf("failed to use shader %s: %w", s.name, err)
		}
		activeShader = s.id
	}
	return nil
}

func (s *glshader) checkState(uniform string, value any) bool {
	current, exists := s.state[uniform]
	if !exists || current != value {
		s.state[uniform] = value
		return false
	}
	return true
}

// SetFragmentData sets the name of the fragment color output variable
func (s *glshader) SetFragmentData(fragVariable string) {
	if err := gl.BindFragDataLocation(s.id, fragVariable); err != nil {
		panic(err)
	}
}

// Attach a shader to the program. Panics if the program is already linked
func (s *glshader) Attach(stage shader.Stage) {
	if s.linked {
		panic(fmt.Sprintf("cant attach, shader %s is already linked", s.name))
	}
	gl.AttachShader(s.id, stage.ID())
	s.shaders = append(s.shaders, stage)
}

// Link the currently attached shaders into a program. Panics on failure
func (s *glshader) Link() error {
	if s.linked {
		return nil
	}

	if err := gl.LinkProgram(s.id); err != nil {
		return fmt.Errorf("failed to compile %s: %s", s.name, err)
	}

	s.readAttributes()
	s.readUniforms()

	s.linked = true
	return nil
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
	if s.checkState(name, mat4) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformMatrix4f(s.id, uniform, mat4)
	} else {
		return err
	}
}

// Vec2 sets a Vec2 uniform value
func (s *glshader) Vec2(name string, vec vec2.T) error {
	if s.checkState(name, vec) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec2f(s.id, uniform, vec)
	} else {
		return err
	}
}

// Vec3 sets a Vec3 uniform value
func (s *glshader) Vec3(name string, vec vec3.T) error {
	if s.checkState(name, vec) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec3f(s.id, uniform, vec)
	} else {
		return err
	}
}

func (s *glshader) Vec3Array(name string, vecs []vec3.T) error {
	// unable to easily compare arrays to check state
	if uniform, err := s.Uniform(fmt.Sprintf("%s[0]", name)); err == nil {
		return gl.UniformVec3fArray(s.id, uniform, vecs)
	} else {
		return err
	}
}

// Vec4 sets a Vec4f uniform value
func (s *glshader) Vec4(name string, vec vec4.T) error {
	if s.checkState(name, vec) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec4f(s.id, uniform, vec)
	} else {
		return err
	}
}

// Int32 sets an integer 32 uniform value
func (s *glshader) Int32(name string, value int) error {
	if s.checkState(name, value) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		if name == "depth" {
			fmt.Println(uniform)
		}
		return gl.UniformVec1i(s.id, uniform, value)
	} else {
		return err
	}
}

// Uint32 sets an unsigned integer 32 uniform value
func (s *glshader) Uint32(name string, value int) error {
	if s.checkState(name, value) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec1ui(s.id, uniform, value)
	} else {
		return err
	}
}

func (s *glshader) Bool(name string, value bool) error {
	if s.checkState(name, value) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformBool(s.id, uniform, value)
	} else {
		return err
	}
}

// Float sets a float uniform value
func (s *glshader) Float(name string, value float32) error {
	if s.checkState(name, value) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec1f(s.id, uniform, value)
	} else {
		return err
	}
}

// RGB sets a uniform to a color RGB value
func (s *glshader) RGB(name string, color color.T) error {
	if s.checkState(name, color) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec3f(s.id, uniform, vec3.T{
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
	if s.checkState(name, color) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformVec4f(s.id, uniform, vec4.T{
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
	if s.checkState(name, slot) {
		return nil
	}
	if uniform, err := s.Uniform(name); err == nil {
		return gl.UniformTexture2D(s.id, uniform, slot)
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
	count := gl.GetActiveAttributeCount(s.id)

	s.attributes = make(shader.AttributeMap, count)
	for i := 0; i < count; i++ {
		attr := gl.GetActiveAttribute(s.id, i)
		s.attributes[attr.Name] = attr
		fmt.Println(s.name, "uniform", attr.Name, attr.Type, "=", attr.Index, "size:", attr.Size)
	}
}

func (s *glshader) readUniforms() {
	count := gl.GetActiveUniformCount(s.id)
	s.uniforms = make(map[string]shader.UniformDesc, count)
	for i := 0; i < count; i++ {
		uniform := gl.GetActiveUniform(s.id, i)
		s.uniforms[uniform.Name] = uniform
		fmt.Println(s.name, "uniform", uniform.Name, uniform.Type, "=", uniform.Index, "size:", uniform.Size)
	}
}
