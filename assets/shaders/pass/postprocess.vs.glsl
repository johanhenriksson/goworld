#version 330

in vec3 position;
in vec2 texcoord_0;

out vec2 texcoord0;

void main() {
    /* pass through texture coordinate */
    texcoord0 = texcoord_0;

    /* set position - coordinates are already in clip space */
    gl_Position = vec4(position, 1.0);
}