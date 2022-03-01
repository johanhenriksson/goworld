#version 330

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

layout(location=0) in vec3 position;
in vec3 normal;
in vec4 color_0;

out vec3 position0;
out vec3 normal0;
out vec4 color0;

void main() {
    // compute modelview matrix
    // perhaps do this offline?
    mat4 mv = view * model;

    // gbuffer view space normal
    normal0 = normalize((mv * vec4(normal, 0.0)).xyz);

    // gbuffer view space position
    position0 = (mv * vec4(position, 1.0)).xyz;

    // pass vertex color
    color0 = color_0;

    // finally, transform view -> clip space and output vertex position
    gl_Position = projection * vec4(position0, 1.0);
}
