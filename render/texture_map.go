package render

type TextureMap struct {
	shader   *Shader
	textures map[string]*Texture
	slots    []string
}

func NewTextureMap(shader *Shader) *TextureMap {
	return &TextureMap{
		shader:   shader,
		slots:    make([]string, 0, 8),
		textures: make(map[string]*Texture, 8),
	}
}

// AddTexture attaches a new texture to this material, and assings it to the next available texture slot.
func (tm *TextureMap) Add(name string, tex *Texture) {
	if _, exists := tm.textures[name]; exists {
		tm.Set(name, tex)
		return
	}
	tm.textures[name] = tex
	tm.slots = append(tm.slots, name)
}

// Set changes a bound texture
func (tm *TextureMap) Set(name string, tex *Texture) {
	if _, exists := tm.textures[name]; !exists {
		panic("no such texture")
	}
	tm.textures[name] = tex
}

func (tm *TextureMap) Length() int {
	return len(tm.slots)
}

func (tm *TextureMap) Get(name string) *Texture {
	if tx, exists := tm.textures[name]; exists {
		return tx
	}
	return nil
}

func (tm *TextureMap) Slot(i int) *Texture {
	if i < 0 || i >= len(tm.slots) {
		return nil
	}
	return tm.textures[tm.slots[i]]
}

func (tm *TextureMap) Use() {
	for i, name := range tm.slots {
		tex := tm.textures[name]
		tex.Use(i)
		tm.shader.Int32(name, i)
	}
}
