#version 450

#include "lib/common.glsl"
#include "lib/ui.glsl"

OUT(0, vec4, color)
OUT(1, vec2, texcoord)
OUT(2, vec2, position)
OUT(3, flat vec2, center)
OUT(4, flat vec2, half_size)
OUT(5, flat uint, quad_index)

out gl_PerVertex 
{
	vec4 gl_Position;   
};

const ivec3 vertices[] =
{
	{-1, -1, 0}, // 0
	{-1, +1, 1}, // 1
	{+1, -1, 2}, // 2

	{-1, +1, 1}, // 1
	{+1, +1, 3}, // 3
	{+1, -1, 2}, // 2
};

void main() 
{
	out_quad_index = gl_InstanceIndex;
	Quad quad = quads.item[out_quad_index];
	vec2 vertex = vertices[gl_VertexIndex].xy;
	int corner = vertices[gl_VertexIndex].z;

	out_half_size = (quad.max - quad.min) / 2;
	out_center = (quad.max + quad.min) / 2;
	out_position = vertex * out_half_size + out_center;

	vec2 tex_half_size = (quad.uv_max - quad.uv_min) / 2;
	vec2 tex_center = (quad.uv_max + quad.uv_min) / 2;
	out_texcoord = vertex * tex_half_size + tex_center;

	gl_Position = vec4(
		2 * out_position.x / config.resolution.x - 1,
		2 * out_position.y / config.resolution.y - 1,
		1 - quad.zindex / (config.zmax + 1),
		1);
	
	out_color = quad.color[corner];
}
