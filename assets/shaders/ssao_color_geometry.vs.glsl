#version 330

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

in vec3 position;
in int normal_id;
in vec3 color;
in float occlusion;

out vec4 color0;
out vec3 normal0;
out vec3 position0;

// normal lookup table
const vec3 normals[7] = vec3[7] (
    vec3(0,0,0),  // normal 0 - undefined
    vec3(1,0,0),  // x+
    vec3(-1,0,0), // x-
    vec3(0,1,0),  // y+
    vec3(0,-1,0), // y-
    vec3(0,0,1),  // z+
    vec3(0,0,-1)  // z-
);

void main() {
    // compute modelview matrix
    // perhaps do this offline?
    mat4 mv = view * model;

    // gbuffer view space normal
    vec3 normal = normals[normal_id];
    normal0 = normalize((mv * vec4(normal, 0.0)).xyz);

    // gbuffer view space position
    position0 = (mv * vec4(position, 1.0)).xyz;

    // pass color and occlusion
    color0 = vec4(color, occlusion);

    // finally, transform view -> clip space and output vertex position
    gl_Position = projection * vec4(position0,1);
}
