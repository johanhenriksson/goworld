#version 330

uniform sampler2D image;
uniform vec4 tint;

in vec2 out_uv;
out vec4 frag_color;

void main() {
    frag_color = texture(image, out_uv) * tint;
}
