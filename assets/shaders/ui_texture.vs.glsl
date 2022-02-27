#version 330

uniform mat4 viewport;
uniform mat4 model;

layout(location=0) in vec3 position;
in vec4 color;
in vec2 texcoord;

out vec2 out_uv;
out vec4 out_color;

void main() {
    out_uv      = texcoord;
    out_color   = color;
    gl_Position = viewport * model * vec4(position, 1);
}
