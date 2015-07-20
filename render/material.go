package render

type MaterialTextureMap map[int32]*Texture

type Material struct {
    Shader      *ShaderProgram
    Textures    MaterialTextureMap
}

func CreateMaterial() *Material {
    return &Material {
        Textures: make(MaterialTextureMap),
    }
}
