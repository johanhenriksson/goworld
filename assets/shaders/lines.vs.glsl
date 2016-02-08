#version 330

uniform mat4 mvp;

in vec3 vertex;
in vec4 color;

out vec4 out_color;

void main() {
    out_color   = color;
    gl_Position = mvp * vec4(vertex,1);
}
