#version 330

const vec3 lightPos = vec3(-2, 2, -1);

uniform bool gbuffer_output = true;

in vec3 position0;
in vec3 normal0;
in vec4 color0;

layout(location=0) out vec4 out_diffuse;
layout(location=1) out vec4 out_normal;
layout(location=2) out vec4 out_position;

void main() {
    vec3 dir = normalize(lightPos);
    float contrib = max(dot(dir, normal0), 0.0);

    out_diffuse = vec4((0.5 + 0.5 * contrib) * color0.rgb, color0.a);

    if (gbuffer_output) {
        // store position in gbuffer
        out_position = vec4(position0, 1.0);

        // pack normal after interpolation
        // store in gbuffer
        out_normal = vec4((normal0 + 1.0) / 2.0, 1);
    }
}
