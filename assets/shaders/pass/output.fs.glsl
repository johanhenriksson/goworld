#version 330

in vec2 texcoord0;
layout(location=0) out vec4 color;

uniform sampler2D tex_input;
uniform sampler2D tex_depth;

void main() {
    color = texture(tex_input, texcoord0);

    // restore depth buffer
    float depth = texture(tex_depth, texcoord0).r;
    gl_FragDepth = depth;
}