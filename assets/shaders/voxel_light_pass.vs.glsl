#version 330

uniform mat4 model;
uniform mat4 viewport;

in vec3 position;
in vec2 texcoord;

out vec2 texcoord0;

void main() {
    texcoord0 = texcoord;
    vec4 p = viewport * model * vec4(1);
    gl_Position = vec4(position, 1.0) + 0.0001 * p;
}
