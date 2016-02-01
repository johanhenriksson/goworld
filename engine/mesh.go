package engine

import (
    "github.com/johanhenriksson/goworld/render"
)

/* Meshes connect vertex data to shaders through the use of Materials */
type Mesh struct {
    VertexData    *render.VertexArray
    Material    *render.Material
}

func CreateMesh(data *render.VertexArray, mat *render.Material) *Mesh {
    mesh := &Mesh {
        VertexData: data,
        Material: mat,
    }
    data.Bind()
    mat.Setup()
    return mesh
}

func (mesh *Mesh) Render() {
    mesh.Material.Use()
    mesh.VertexData.Bind()
    mesh.VertexData.Draw()
}
