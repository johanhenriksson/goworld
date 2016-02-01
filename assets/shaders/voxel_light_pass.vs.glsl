#version 330

// legacy. to be removed asap
uniform mat4 model;
uniform mat4 viewport;

in vec3 position;
in vec2 texcoord;

out vec2 texcoord0;

void main() {
    /* pass through texture coordinate */
    texcoord0 = texcoord;

    // ugly hack to work with the old stuff...
    vec4 p = viewport * model * vec4(1);

    /* set position - coordinates are already in clip space */
    gl_Position = vec4(position, 1.0) + 0.0001 * p;
}
