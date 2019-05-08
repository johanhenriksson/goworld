#version 330

const float TileSize = 16;
const float TilesetTexWidth = 4096;
const float TilesetTexHeight = 4096;

uniform mat4 model;
uniform mat4 mvp;

in vec3 position;
in vec3 normal;
in vec2 tile;

out vec2 texcoord0;
out vec3 normal0;
out vec3 position0;

void main() {
    /* Transform normal */
    normal0 = normalize((model * vec4(normal,0)).xyz);

    /* Convert tile coordinate to texture coord */
    texcoord0 = vec2(tile.x, 1.0 - tile.y) * TileSize / TilesetTexHeight;

    gl_Position = mvp * vec4(position, 1);
    position0 = gl_Position;
}
