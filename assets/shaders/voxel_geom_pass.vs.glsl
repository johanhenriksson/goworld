#version 330

const float TileSize = 16;
const float TilesetTexWidth = 4096;
const float TilesetTexHeight = 4096;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 position;
in vec3 normal;
in vec2 tile;

out vec2 texcoord0;
out vec3 normal0;

void main() {

    /* Transform normal */
    normal0 = normalize((model * vec4(normal,1)).xyz);

    /* Convert tile coordinate to texture coord */
    texcoord0 = vec2(tile.x, 1.0 - tile.y) * TileSize / TilesetTexHeight;

    gl_Position = projection * camera * model * vec4(position,1);
}
