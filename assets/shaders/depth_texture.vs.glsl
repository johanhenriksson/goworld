#version 330

uniform mat4 viewport;
uniform mat4 model;

in vec3 vertex;
in vec2 uv;

out vec2 out_uv;

void main() {
    out_uv      = uv;
    gl_Position = viewport * model * vec4(vertex,1);
}
