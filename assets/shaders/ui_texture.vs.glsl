#version 330

uniform mat4 viewport;
uniform mat4 model;

layout(location=0) in vec3 position;
in vec4 color_0;
in vec2 texcoord_0;

out vec2 out_uv;
out vec4 out_color;

void main() {
    out_uv      = texcoord_0;
    out_color   = color_0;
    gl_Position = viewport * model * vec4(position, 1);
}
