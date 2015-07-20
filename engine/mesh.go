package engine

import (
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/geometry"
)

type Mesh struct {
    Geometry    *geometry.VertexArray
    Material    *render.Material
}

func CreateMesh(data *geometry.VertexArray, mat *render.Material) *Mesh {
    mesh := &Mesh {
        Geometry: data,
        Material: mat,
    }
    data.Bind()
    mat.Setup()
    return mesh
}

func (mesh *Mesh) Render() {
    mesh.Material.Use()
    mesh.Geometry.Bind()
    mesh.Geometry.Draw()
}
