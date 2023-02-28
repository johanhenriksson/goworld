#version 450

#extension GL_ARB_separate_shader_objects : enable
#extension GL_ARB_shading_language_420pack : enable

// Attributes
layout (location = 0) in vec3 position;

// Uniforms

struct Quad {
    vec2 min; // top left
    vec2 max; // bottom right
    vec2 uv_min; // top left uv
    vec2 uv_max; // bottom right uv
    vec4 color;
    uint texture;
    float zindex;
};

layout (binding = 0) uniform Config {
    vec2 resolution;
} config;


layout (binding = 1) readonly buffer QuadBuffer {
	Quad quads[];
} ssbo;

// Varyings
layout (location = 0) out vec4 color0;
layout (location = 1) flat out uint texture0;
layout (location = 2) out vec2 uv0;

out gl_PerVertex 
{
    vec4 gl_Position;   
};

const vec2 vertices[] =
{
    {-1, -1},
    {-1, +1},
    {+1, -1},
    {+1, +1},
};

void main() 
{
	Quad quad = ssbo.quads[gl_InstanceIndex];

    vec2 dst_half_size = (quad.max - quad.min) / 2;
    vec2 dst_center = (quad.max + quad.min) / 2;
    vec2 dst_pos = vertices[gl_VertexIndex] * dst_half_size + dst_center;

    vec2 tex_half_size = (quad.uv_max - quad.uv_min) / 2;
    vec2 tex_center = (quad.uv_max + quad.uv_min) / 2;
    vec2 tex_pos = vertices[gl_VertexIndex] * tex_half_size + tex_center;

    gl_Position = vec4(
        2 * dst_pos.x / config.resolution.x - 1,
        2 * dst_pos.y / config.resolution.y - 1,
        0,
        1);
    
    uv0 = tex_pos;
	color0 = quad.color;
    texture0 = quad.texture;
}