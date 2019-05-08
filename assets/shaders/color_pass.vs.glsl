#version 330

in vec3 position;
in vec2 texcoord;

out vec2 texcoord0;

void main() {
    texcoord0 = texcoord;
    gl_Position = vec4(position, 1);
}