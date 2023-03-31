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
	vec4 color[4];
	float zindex;
	float corner_radius;
	float edge_softness;
	float border;
	uint texture;
};

layout (binding = 0) uniform Config {
	vec2 resolution;
	float zmax;
} config;


layout (binding = 1) readonly buffer QuadBuffer {
	Quad quads[];
} ssbo;

// Varyings
layout (location = 0) out vec4 color0;
layout (location = 1) out vec2 uv0;
layout (location = 2) out vec2 pos;
layout (location = 3) flat out uint texture0;
layout (location = 4) flat out vec2 center;
layout (location = 5) flat out vec2 half_size;
layout (location = 6) flat out float corner_radius;
layout (location = 7) flat out float edge_softness;
layout (location = 8) flat out float border;


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

	half_size = (quad.max - quad.min) / 2;
	center = (quad.max + quad.min) / 2;
	pos = vertices[gl_VertexIndex] * half_size + center;

	vec2 tex_half_size = (quad.uv_max - quad.uv_min) / 2;
	vec2 tex_center = (quad.uv_max + quad.uv_min) / 2;
	uv0 = vertices[gl_VertexIndex] * tex_half_size + tex_center;

	gl_Position = vec4(
		2 * pos.x / config.resolution.x - 1,
		2 * pos.y / config.resolution.y - 1,
		1 - quad.zindex / (config.zmax+1),
		1);
	
	color0 = quad.color[gl_VertexIndex];
	texture0 = quad.texture;
	corner_radius = quad.corner_radius;
	edge_softness = quad.edge_softness;
	border = quad.border;
}