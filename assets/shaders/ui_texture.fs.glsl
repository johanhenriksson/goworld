#version 330

uniform sampler2D tex0;

in vec2 out_uv;
out vec4 frag_color;

void main() {
    frag_color = texture(tex0, out_uv);
}
