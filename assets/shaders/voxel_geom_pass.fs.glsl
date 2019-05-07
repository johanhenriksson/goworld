#version 330

uniform sampler2D tileset;

in vec2 texcoord0;
in vec3 normal0;
in vec3 position0;

layout(location=0) out vec4 out_diffuse;
layout(location=1) out vec4 out_normal;
layout(location=2) out vec3 out_position;

void main() {
    out_diffuse = vec4(texture(tileset, texcoord0).rgb, 0.0);

    vec4 pack_normal = vec4((normal0 + 1.0) / 2.0, 1);
    out_normal = pack_normal;
}
