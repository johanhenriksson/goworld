#version 330

uniform mat4 mvp;

/* input */
in vec3 position;

void main() {
    gl_Position = mvp * vec4(position,1);
}
