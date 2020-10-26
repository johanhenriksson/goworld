#version 330

uniform sampler2D image;
uniform vec4 tint;
uniform bool invert = false;
uniform bool depth = false;

in vec2 out_uv;
out vec4 frag_color;

void main() {
    vec2 uv = out_uv;
    if (invert) {
        uv.y = 1 - uv.y;
    }

    vec4 color = texture(image, uv);
    if (depth) {
        color = vec4(vec3(color.r), 1);
    }


    frag_color = color * tint;
}
