#version 330

uniform mat4 viewport;
uniform mat4 model;

in vec3 position;
in vec2 texcoord;

out vec2 out_uv;

void main() {
    out_uv      = texcoord;
    gl_Position = viewport * model * vec4(position, 1);
}
