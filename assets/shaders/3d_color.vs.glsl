#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vertex;
in vec3 color;

out vec3 fragColor;

void main() {
    fragColor = color;
    gl_Position = projection * camera * model * vec4(vertex, 1);
}
