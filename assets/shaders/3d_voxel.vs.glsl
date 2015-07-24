#version 330

const float TileSize = 16;
const float TilesetTexWidth = 4096;
const float TilesetTexHeight = 4096;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

uniform vec3 cameraPos;
uniform vec3 lightPos;

in vec3 vertex;
in vec3 normal;
in vec2 tile;

out vec2 texcoord;
out vec3 worldNormal;
out vec3 L;
out vec3 V;
out float lightDistance;

void main() {
    /* Transform normal */
    mat3 normalMatrix = mat3(model);
    worldNormal = normalMatrix * normal;

    /* Convert tile coordinate to texture coord */
    texcoord = vec2(tile.x, 1.0 - tile.y) * TileSize / TilesetTexHeight;

    /* Vertex position */
    vec4 worldPos = model * vec4(vertex, 1);

    /* Lighting */
    L = lightPos - worldPos.xyz;
    lightDistance = length(L);
    L = normalize(L);
    V = normalize(cameraPos - worldPos.xyz);

    gl_Position = projection * camera * worldPos;
}
