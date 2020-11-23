#version 330

in vec3 position0;
in vec3 normal0;
in vec4 color0;

layout(location=0) out vec4 out_diffuse;
layout(location=1) out vec4 out_normal;
layout(location=2) out vec4 out_position;

void main() {
    // store color in gbuffer
    out_diffuse = vec4(color0.rgb, 1);

    // pack normal after interpolation
    // store in gbuffer
    vec4 pack_normal = vec4((normal0 + 1.0) / 2.0, 1);
    out_normal = pack_normal;

    // store position in gbuffer
    out_position = vec4(position0, 1.0);
}
