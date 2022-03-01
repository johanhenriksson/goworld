#version 330

uniform mat4 mvp;

layout(location=0) in vec3 position;
in vec4 color_0;

out vec4 out_color;

void main() {
    out_color   = color_0;
    gl_Position = mvp * vec4(position, 1);
}
