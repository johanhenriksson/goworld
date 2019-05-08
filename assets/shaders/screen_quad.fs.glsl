#version 330

in vec2 texcoord0;
layout(location=0) out vec4 color;

uniform sampler2D tex_input;

void main() {
    color = texture(tex_input, texcoord0);
}