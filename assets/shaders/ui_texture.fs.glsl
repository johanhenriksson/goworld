#version 330

uniform sampler2D image;

in vec2 out_uv;
in vec4 out_color;
out vec4 frag_color;

void main() {
    vec2 uv = out_uv;

    // uv's outside 0.0-1.0 should be transparent/discarded
    if (all(lessThan(uv, vec2(0.0))) && all(greaterThan(uv, vec2(1.0)))) {
        discard;
    }

    vec4 tint = vec4(1);
    if (out_color.a > 0) {
        tint = out_color;
    }

    vec4 texcolor = texture(image, uv);
    frag_color = texcolor.rgba * tint;
}
