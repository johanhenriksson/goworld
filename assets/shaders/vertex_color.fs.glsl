#version 330

const vec3 lightPos = vec3(-2, 2, -1);

in vec3 normal0;
in vec4 color0;

out vec4 frag_color;

void main() {
    vec3 dir = normalize(lightPos);
    float contrib = max(dot(dir, normal0), 0.0);

    frag_color = vec4((0.5 + 0.5 * contrib) * color0.rgb, color0.a);
}
