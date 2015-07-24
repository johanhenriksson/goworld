#version 330

const float TileSize = 16;
const float TilesetTexWidth = 4096;
const float TilesetTexHeight = 4096;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vertex;
in vec3 normal;
in vec2 tile;

out vec2 texcoord;
out vec3 frag_normal;

void main() {
    /* Transform normal */
    mat3 normalMatrix = mat3(model);
    frag_normal = normalMatrix * normal;

    /* Convert tile coordinate to texture coord */
    texcoord = vec2(tile.x, 1.0 - tile.y) * TileSize / TilesetTexHeight;

    gl_Position = projection * camera * model * vec4(vertex, 1);
}
