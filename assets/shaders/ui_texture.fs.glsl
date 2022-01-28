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

    // depth map/single channel/grayscale images
    // todo: rename this parameter
    if (depth) {
        color = vec4(vec3(color.r), 1);
    }

    // uv's outside 0.0-1.0 should be transparent/discarded
    if (all(lessThan(uv, vec2(0.0))) && all(greaterThan(uv, vec2(1.0)))) {
        discard;
    }

    frag_color = color * tint;
}
