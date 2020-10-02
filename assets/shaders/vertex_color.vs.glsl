#version 330

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

in vec3 position;
in vec3 normal;
in vec4 color;

out vec3 position0;
out vec3 normal0;
out vec4 color0;

void main() {
    color0 = color;
    mat4 mp = projection * view * model;
    normal0 = normalize((model * vec4(normal, 0.0)).xyz);
    gl_Position = mp * vec4(position, 1);
    position0 = gl_Position.xyz;
}
