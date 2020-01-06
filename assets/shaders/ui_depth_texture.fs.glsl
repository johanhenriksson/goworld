#version 330

uniform sampler2D image;

in vec2 out_uv;
out vec4 frag_color;

void main() {
    float depth = texture(image, out_uv).r;
    frag_color = vec4(depth,depth,depth,1);
}
