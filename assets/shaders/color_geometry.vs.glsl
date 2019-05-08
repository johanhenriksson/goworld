#version 330

uniform mat4 model;
uniform mat4 mvp;

in vec3 position;
in int normal_id;
in vec3 color;
in float occlusion;

out vec4 color0;
out vec3 normal0;
out vec3 position0;

/* In order to save vertex memory, normals are looked up in a table */
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
    /* Transform normal */
    vec3 normal = normals[normal_id];
    normal0 = normalize((model * vec4(normal,0)).xyz);

    /* pass color */
    color0 = vec4(color, occlusion);

    gl_Position = mvp * vec4(position,1);
    position0 = (model * vec4(position, 1.0)).xyz;
}
