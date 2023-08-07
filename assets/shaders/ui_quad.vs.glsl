#version 450
#extension GL_GOOGLE_include_directive : enable

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

const vec2 vertices[] =
{
	{-1, -1},
	{-1, +1},
	{+1, -1},
	{+1, +1},
};

void main() 
{
	out_quad_index = gl_InstanceIndex;
	Quad quad = quads.item[out_quad_index];

	out_half_size = (quad.max - quad.min) / 2;
	out_center = (quad.max + quad.min) / 2;
	out_position = vertices[gl_VertexIndex] * out_half_size + out_center;

	vec2 tex_half_size = (quad.uv_max - quad.uv_min) / 2;
	vec2 tex_center = (quad.uv_max + quad.uv_min) / 2;
	out_texcoord = vertices[gl_VertexIndex] * tex_half_size + tex_center;

	gl_Position = vec4(
		2 * out_position.x / config.resolution.x - 1,
		2 * out_position.y / config.resolution.y - 1,
		1 - quad.zindex / (config.zmax + 1),
		1);
	
	out_color = quad.color[gl_VertexIndex];
}
