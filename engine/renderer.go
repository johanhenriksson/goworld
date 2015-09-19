package engine

import (
    "github.com/go-gl/gl/v4.1-core/gl"
    mgl "github.com/go-gl/mathgl/mgl32"
    "github.com/johanhenriksson/goworld/render"
)

type Component interface {
    Update(float32)
    Draw(render.DrawArgs)
}

type Renderer struct {
    Width       int32
    Height      int32
    Geometry    *render.GeometryBuffer
    Scene       *Scene
    /* TODO output quad */
    gs          *render.ShaderProgram
}

func NewRenderer(width, height int32) *Renderer {
    r := &Renderer {
        Width: width,
        Height: height,
        Geometry: render.CreateGeometryBuffer(width, height),
        Scene: NewScene(),
        gs: render.CompileVFShader("/assets/shaders/voxel_geom"),
    }
    r.Geometry.Unbind()
    program := r.gs
    program.Use()
    program.Vec3("lightPos", &mgl.Vec3{ 8,15,8 })
    program.Float("lightIntensity", 250.0)
    program.Float("ambient", 0.6)
    m := mgl.Ident4()
    program.Matrix4f("model", &m[0])
    return r
}

func (r *Renderer) Draw() {
    r.Geometry.Bind()
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
    cam := r.Scene.Camera
    r.gs.Use()
    r.gs.Matrix4f("projection", &cam.Projection[0])
    r.gs.Matrix4f("camera", &cam.View[0])
    r.gs.Vec3("cameraPos", &cam.Transform.Position)
    r.Scene.Draw(r.gs)
    r.Geometry.Unbind()
}

type Scene struct {
    Camera      *Camera
    Objects     []*Object
}

func NewScene() *Scene {
    return &Scene {
        Objects: []*Object { },
    }
}

func (s *Scene) Add(object *Object) {
    /* TODO look for lights */
    s.Objects = append(s.Objects, object)
}

func (s *Scene) Draw(shader *render.ShaderProgram) {
    if s.Camera == nil {
        return
    }
    args := render.DrawArgs {
        Viewport: s.Camera.View,
        Transform: mgl.Ident4(),
        Shader: shader,
    }
    for _, obj := range s.Objects {
        obj.Draw(args)
    }
}

func (s *Scene) Update(dt float32) {
    for _, obj := range s.Objects {
        obj.Update(dt)
    }
}


/** Game object */
type Object struct {
    *Transform
    Components  []Component
    Children    []*Object
}

func NewObject(x,y,z float32) *Object {
    return &Object {
        Transform: CreateTransform(x,y,z),
        Components: []Component { },
        Children: []*Object { },
    }
}

func (o *Object) Attach(component Component) {
    o.Components = append(o.Components, component)
}

func (o *Object) Draw(args render.DrawArgs) {
    /* Apply transform */
    args.Transform = o.Transform.Matrix.Mul4(args.Transform)
    args.Shader.Matrix4f("model", &args.Transform[0])
    for _, comp := range o.Components { comp.Draw(args) }
    for _, child := range o.Children { child.Draw(args) }
}

func (o *Object) Update(dt float32) {
    for _, comp := range o.Components { comp.Update(dt) }
    for _, child := range o.Children { child.Update(dt) }
}

type Light struct {
    /* TODO types: point, spot, directional */
    Intensity float32
    Range float32
}

func (l *Light) Update(dt float32) {
}
func (l *Light) Draw(args render.DrawArgs) {
}

